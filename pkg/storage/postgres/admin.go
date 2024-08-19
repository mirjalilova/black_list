package postgres

import (
	"database/sql"

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

	res.Count = int32(len(res.Hr))

	return &res, nil
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

	return res, nil
}
