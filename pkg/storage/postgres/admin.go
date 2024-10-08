package postgres

import (
	"database/sql"
	"fmt"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (s *AdminRepo) Approve(req *pb.CreateHR) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
	}

	query := `
			INSERT INTO hr 
				(user_id, approved_by) 
			VALUES 
				($1, $2)`

	_, err = tr.Exec(query, req.UserId, req.ApprovedBy)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	query = `UPDATE 
				users 
			SET 
				role = 'hr',
				updated_at = now()
			WHERE 
				id = $1 AND deleted_at = 0`

	_, err = tr.Exec(query, req.UserId)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	tr.Commit()

	return res, nil
}

func (s *AdminRepo) ListHR(req *pb.Filter) (*pb.GetAllHRRes, error) {
	var res pb.GetAllHRRes

	query := `SELECT 
				h.id,
				u.full_name,
				u.email,
				u.date_of_birth,
				h.created_at
			FROM hr h 
			JOIN users u ON h.user_id = u.id
			WHERE h.deleted_at=0 LIMIT $1 OFFSET $2`

	req.Offset = (req.Offset - 1) * req.Limit

	rows, err := s.db.Query(query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hr pb.Hr

		err := rows.Scan(
			&hr.Id,
			&hr.FullName,
			&hr.Email,
			&hr.DateOfBirth,
			&hr.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		res.Hr = append(res.Hr, &hr)
	}

	query = `SELECT COUNT(*) FROM hr WHERE deleted_at=0`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)

	return &res, nil
}

func (s *AdminRepo) GetHRById(req *pb.GetById) (*pb.Hr, error) {
	res := &pb.Hr{}

	query := `SELECT
				h.id,
				u.full_name,
				u.email,
				u.date_of_birth,
				h.created_at
			FROM
				hr h 
            JOIN 
				users u ON h.user_id = u.id
			WHERE
				h.id = $1 AND h.deleted_at = 0`
	
	err := s.db.QueryRow(query, req.Id).Scan(
		    &res.Id,
            &res.FullName,
            &res.Email,
            &res.DateOfBirth,
            &res.CreatedAt,
        )
    if err!= nil {
		return nil, err
	}

	return res, nil
}

func (s *AdminRepo) Delete(req *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

	tr, err := s.db.Begin()
	if err != nil {
		return res, err
	}

	query := `
			UPDATE 
				hr 
			SET 	
				deleted_at = EXTRACT(EPOCH FROM NOW()) 
			WHERE 
				id = $1`

	_, err = tr.Exec(query, req.Id)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	var user_id string
	query = `SELECT user_id FROM hr WHERE id = $1`
	err = tr.QueryRow(query, req.Id).Scan(&user_id)
	if err != nil {
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
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	tr.Commit()

	return res, nil
}

func (s *AdminRepo) GetAllUsers(req *pb.ListUserReq) (*pb.ListUserRes, error) {
	res := &pb.ListUserRes{}

	query := `SELECT 
				id, 
				username, 
				full_name,
				email, 
				date_of_birth,
				role 
			FROM 
				users 
			WHERE 
				deleted_at=0 AND role = 'user'`

	var args []interface{}

	if req.Username != "" && req.Username != "string" {
		args = append(args, "%"+req.Username+"%")
		query += fmt.Sprintf(" AND username ILIKE $%d", len(args))
	}

	if req.FullName != "" && req.FullName != "string" {
		args = append(args, "%"+req.FullName+"%")
		query += fmt.Sprintf(" AND full_name ILIKE $%d", len(args))
	}

	req.Filter.Offset = (req.Filter.Offset - 1) * req.Filter.Limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, req.Filter.Limit, req.Filter.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user pb.UserRes
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.FullName,
			&user.Email,
			&user.DateOfBirth,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}
		res.Users = append(res.Users, &user)
	}

	query = `SELECT COUNT(*) FROM users WHERE deleted_at=0 AND role = 'user'`
	var count int64
	err = s.db.QueryRow(query).Scan(&count)
	if err!= nil {
        return nil, err
    }

	res.Count = int32(count)

	return res, nil
}

func (s *AdminRepo) ChangeRole(req *pb.ChangeRoleReq) (*pb.Void, error) {
	res := &pb.Void{}

	tr , err := s.db.Begin()
	if err!= nil {
		return res, err
    }

	query := `UPDATE 
				users
			SET
			    role = $1
			WHERE
				id = $2 AND deleted_at = 0`

	_, err = tr.Exec(query, req.Role, req.UserId)
	if err!= nil {
        return res, err
    }

	if req.Role == "admin" {
		query = `INSERT INTO hr (user_id, approved_by) VALUES ($1, $2)`
		_, err = tr.Exec(query, req.UserId, req.UserId)
		if err!= nil {
			tr.Rollback()
			return nil, err
		}
	}

	tr.Commit()

	return res, nil
}

func (s *AdminRepo) GetUserById(req *pb.GetById) (*pb.UserRes, error) {
	res := &pb.UserRes{}

	query := `SELECT 
                id, 
                username, 
                full_name,
                email, 
                date_of_birth,
                role 
            FROM 
                users 
            WHERE 
                id = $1 AND deleted_at=0`

	err := s.db.QueryRow(query, req.Id).Scan(
		&res.Id,
		&res.Username,
		&res.FullName,
		&res.Email,
		&res.DateOfBirth,
		&res.Role,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return res, nil
}