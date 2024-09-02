package postgres

import (
	"database/sql"
	"fmt"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type BlackListRepo struct {
	db *sql.DB
}

func NewBlackListRepo(db *sql.DB) *BlackListRepo {
	return &BlackListRepo{
		db: db,
	}
}

func (s *BlackListRepo) Add(req *pb.BlackListCreate) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
	}

	query := `SELECT id FROM hr WHERE user_id = $1`

	var hr_id string
	err = tr.QueryRow(query, req.AddedBy).Scan(&hr_id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("HR not found for user_id: %s", req.AddedBy)
	} else if err != nil {
		return nil, err
	}

	query = `INSERT INTO black_list 
				(employee_id, reason) 
			VALUES 
				($1, $2)`

	_, err = tr.Exec(query, req.EmployeeId, req.Reason)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `INSERT INTO audit_logs 
				(added_by, employee_id, action)
			VALUES
				($1, $2, 'added')`

	_, err = tr.Exec(query, hr_id, req.EmployeeId)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `UPDATE employees SET is_blocked = true WHERE id = $1`
	_, err = tr.Exec(query, req.EmployeeId)
	if err!= nil {
        tr.Rollback()
        return nil, err
    }

	tr.Commit()

	return res, nil
}

func (s *BlackListRepo) GetAll(req *pb.Filter) (*pb.GetAllBlackListRes, error) {
	res := &pb.GetAllBlackListRes{}

	query := `SELECT 
				u.full_name,
				u.date_of_birth,
				e.position,
				b.reason,
                b.blacklisted_at
			FROM
				users u
			JOIN
				employees e on u.id = e.user_id
			JOIN 
				black_list b on e.id = b.employee_id
			LIMIT $1 OFFSET $2`

	req.Offset = (req.Offset - 1) * req.Limit

	rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		bk := &pb.BlackListRes{}

		err = rows.Scan(
			&bk.FullName,
			&bk.DateOfBirth,
			&bk.Position,
			&bk.Reason,
			&bk.BlacklistedAt,
		)
		if err != nil {
			return nil, err
		}

		res.BlackLists = append(res.BlackLists, bk)
	}

	query = `SELECT COUNT(*) FROM black_list`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)

	return res, nil
}

func (s *BlackListRepo) Remove(req *pb.RemoveReq) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
	}


	query := `SELECT id FROM hr WHERE user_id = $1`

	var hr_id string
	err = tr.QueryRow(query, req.AddedBy).Scan(&hr_id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("HR not found for user_id: %s", req.AddedBy)
	} else if err != nil {
		return nil, err
	}

	query = `DELETE FROM black_list 
				WHERE 
			employee_id = $1`

	_, err = tr.Exec(query, req.EmployeeId)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `INSERT INTO audit_logs
				(added_by, employee_id, action)
			VALUES
				($1, $2, 'removed')`

	_, err = tr.Exec(query, hr_id, req.EmployeeId)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `UPDATE employees SET is_blocked = false WHERE id = $1`
	_, err = tr.Exec(query, req.EmployeeId)
	if err!= nil {
        tr.Rollback()
        return nil, err
    }

	tr.Commit()

	return res, nil
}

func (s *BlackListRepo) MonitoringDailyReport(req *pb.Filter) (*pb.Reports, error) {
	res := &pb.Reports{}

    query := `SELECT 
				u.full_name,
				b.timestamp
            FROM
                audit_logs b 
            JOIN
                employees e on e.id = b.employee_id
			JOIN 
				users u on u.id = e.user_id
			WHERE 
				b.action = 'added' AND b.timestamp >= NOW() - INTERVAL '1 day' LIMIT $1 OFFSET $2`

	req.Offset = (req.Offset - 1) * req.Limit

	rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()

	for rows.Next() {
		report := &pb.Report{}
		err = rows.Scan(
            &report.FullName,
            &report.BlacklistedAt,
        )
		if err!= nil {
            return nil, err
        }

		res.Reports = append(res.Reports, report)
	}

	query = `SELECT COUNT(*) FROM black_list WHERE blacklisted_at >= NOW() - INTERVAL '1 day'`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)
	
	fmt.Println("ssssssss", res)

    return res, nil
}

func (s *BlackListRepo) MonitoringWeeklyReport(req *pb.Filter) (*pb.Reports, error) {
	res := &pb.Reports{}

    query := `SELECT 
				u.full_name,
				b.timestamp
            FROM
                audit_logs b 
            JOIN
                employees e on e.id = b.employee_id
			JOIN 
				users u on u.id = e.user_id
			WHERE 
				b.action = 'added' AND b.timestamp >= NOW() - INTERVAL '1 week' LIMIT $1 OFFSET $2`

	req.Offset = (req.Offset - 1) * req.Limit

	rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()

	for rows.Next() {
		report := &pb.Report{}
		err = rows.Scan(
            &report.FullName,
            &report.BlacklistedAt,
        )
		if err!= nil {
            return nil, err
        }

		res.Reports = append(res.Reports, report)
	}
	
	query = `SELECT COUNT(*) FROM black_list WHERE blacklisted_at >= NOW() - INTERVAL '1 week'`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)	
	
	fmt.Println("ddddd", res)

    return res, nil
}

func (s *BlackListRepo) MonitoringMonthlyReport(req *pb.Filter) (*pb.Reports, error) {
	res := &pb.Reports{}

    query := `SELECT 
				u.full_name,
				b.timestamp
            FROM
                audit_logs b 
            JOIN
                employees e on e.id = b.employee_id
			JOIN 
				users u on u.id = e.user_id
			WHERE 
				b.action = 'added' AND b.timestamp >= NOW() - INTERVAL '1 month' LIMIT $1 OFFSET $2`

	req.Offset = (req.Offset - 1) * req.Limit

	rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()

	for rows.Next() {
		report := &pb.Report{}
		err = rows.Scan(
            &report.FullName,
            &report.BlacklistedAt,
        )
		if err!= nil {
            return nil, err
        }

		res.Reports = append(res.Reports, report)
	}
	
	query = `SELECT COUNT(*) FROM black_list WHERE blacklisted_at >= NOW() - INTERVAL '1 month'`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)		

	fmt.Println("ffffff", res)
    return res, nil
}
