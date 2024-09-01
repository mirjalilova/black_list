package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type HRRepo struct {
	db *sql.DB
}

func NewHRRepo(db *sql.DB) *HRRepo {
	return &HRRepo{
		db: db,
	}
}

func (s *HRRepo) Create(req *pb.EmployeeCreate) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
	}

	query := `SELECT id FROM hr WHERE user_id = $1`

	var hr_id string
	err = tr.QueryRow(query, req.HrId).Scan(&hr_id)
	if err == sql.ErrNoRows {
		tr.Rollback()
		return nil, fmt.Errorf("HR not found for user_id: %s", req.HrId)
	} else if err != nil {
	return nil, err
    }

	query = `INSERT INTO employees 
				(user_id, position, hr_id) 
			VALUES ($1, $2, $3)`

	_, err = tr.Exec(query, req.UserId, req.Position, hr_id)
	if err != nil {
		tr.Rollback()
		return res, err
	}

	query = `UPDATE 
				users
			SET
				role = 'employee',
				updated_at = now()
			WHERE 
				id = $1 AND deleted_at = 0`

	_, err = tr.Exec(query, req.UserId)
	if err != nil {
		tr.Rollback()
		return res, err
	}

	tr.Commit()

	return res, nil
}

func (s *HRRepo) Get(req *pb.GetById) (*pb.Employee, error) {
	res := &pb.Employee{}

	query := `SELECT
				e.id,
				u.full_name,
				u.email,
                u.date_of_birth,
                e.position,
                e.hr_id,
				e.is_blocked
			FROM
				employees e
			JOIN 
				users u ON e.user_id = u.id
			WHERE 
				e.id = $1 AND e.deleted_at = 0`

	err := s.db.QueryRow(query, req.Id).Scan(
		&res.Id,
		&res.FullName,
		&res.Email,
		&res.DateOfBirth,
		&res.Position,
		&res.HrId,
		&res.IsBlocked,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *HRRepo) GetAll(req *pb.ListEmployeeReq) (*pb.ListEmployeeRes, error) {
	res := &pb.ListEmployeeRes{}

	query := `SELECT
				e.id,
				u.full_name,
				u.email,
                u.date_of_birth,
                e.position,
                e.hr_id,
				e.is_blocked
			FROM
				employees e
			JOIN 
				users u ON e.user_id = u.id
			WHERE 
				e.deleted_at = 0`

	var args []interface{}

	if req.Position != "" && req.Position != "string" {
		query += " AND e.position ILIKE $1"
		args = append(args, "%"+req.Position+"%")
	}

	req.Filter.Offset = (req.Filter.Offset - 1) * req.Filter.Limit
	
	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emp pb.Employee
		err := rows.Scan(
			&emp.Id,
			&emp.FullName,
			&emp.Email,
			&emp.DateOfBirth,
			&emp.Position,
			&emp.HrId,
			&emp.IsBlocked,
		)
		if err != nil {
			return nil, err
		}
		res.Employees = append(res.Employees, &emp)
	}

	query = `SELECT COUNT(*) FROM employees WHERE deleted_at=0`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)

	return res, nil
}


func (s *HRRepo) Update(req *pb.UpdateReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE employees SET`

	var arg []interface{}
	var conditions []string

	if req.Position != "" && req.Position != "string" {
		arg = append(arg, req.Position)
		conditions = append(conditions, fmt.Sprintf(" position = $%d", len(arg)))
	}

	if req.HrId != "" && req.HrId != "string" {
		arg = append(arg, req.HrId)
		conditions = append(conditions, fmt.Sprintf(" hr_id = $%d", len(arg)))
	}

	arg = append(arg, req.Id)
	query += strings.Join(conditions, ", ") + fmt.Sprintf(" WHERE id = $%d", len(arg))

	_, err := s.db.Exec(query, arg...)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *HRRepo) Delete(req *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

    tr, err := s.db.Begin()
    if err!= nil {
        return res, err
    }

    query := `UPDATE 
                employees
            SET
                deleted_at = EXTRACT(EPOCH FROM NOW()) 
            WHERE 
                id = $1 AND deleted_at = 0`

    _, err = tr.Exec(query, req.Id)
    if err!= nil {
        return res, err
    }

	var user_id string
	query = `SELECT 
				user_id 
			FROM 
				employees 
			WHERE 
				id = $1`
				
	err = tr.QueryRow(query, req.Id).Scan(&user_id)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

    query = `UPDATE 
                users
            SET
                role = 'user',
                updated_at = now()
            WHERE 
                id = $1 AND deleted_at = 0`

    _, err = tr.Exec(query, user_id)
    if err!= nil {
        return res, err
    }

	tr.Commit()
	
    return res, nil
}