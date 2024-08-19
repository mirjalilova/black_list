package postgres

import (
	"database/sql"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type BalckListRepo struct {
	db *sql.DB
}

func NewBalckListRepo(db *sql.DB) *BalckListRepo {
	return &BalckListRepo{
		db: db,
	}
}

func (s *BalckListRepo) Add(req *pb.BlackListCreate) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
    }

	query := `INSERT INTO black_list 
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
	
	_, err = tr.Exec(query, req.AddedBy, req.EmployeeId)
	if err!= nil {
        tr.Rollback()
        return nil, err
    }

	tr.Commit()

	return res, nil
}

func (s *BalckListRepo) GetAll(req *pb.Filter) (*pb.GetAllBlackListRes, error) {
	res := &pb.GetAllBlackListRes{}

	query := `SELECT 
				u.full_name,
				u.date_of_birth,
				e.position,
				b.reason,
                b.created_at
			FROM
				user u
			JOIN
				employees e on u.id = e.user_id
			JOIN 
				black_list b on e.id = b.employee_id
				LIMIT $1 OFFSET $2`
	
    rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err!= nil {
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
		if err!= nil {
            return nil, err
        }

		res.BlackLists = append(res.BlackLists, bk)
	}

	res.Count = int32(len(res.BlackLists))

	return res, nil
}

func (s *BalckListRepo) Remove(req *pb.RemoveReq) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
    }

	query := `DELETE FROM black_list 
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
	
	_, err = tr.Exec(query, req.AddedBy, req.EmployeeId)
	if err!= nil {
        tr.Rollback()
        return nil, err
    }

	tr.Commit()

	return res, nil
}
