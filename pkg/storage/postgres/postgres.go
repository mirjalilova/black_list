package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mirjalilova/black_list/internal/config"
	"github.com/mirjalilova/black_list/pkg/storage"
)

type Storage struct {
	db         *sql.DB
	AdminS     storage.AdminI
	HRS        storage.HRI
	BlackListS storage.BlackListI
}

func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s sslmode=disable",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	
	return &Storage{
		db:         db,
		AdminS:     NewAdminRepo(db),
		HRS:        NewHRRepo(db),
		BlackListS: NewBalckListRepo(db),
	}, nil
}
func (s *Storage) Admin() storage.AdminI {
	if s.AdminS == nil {
		s.AdminS = NewAdminRepo(s.db)
	}
	return s.AdminS
}

func (s *Storage) HR() storage.HRI {
	if s.HRS == nil {
		s.HRS = NewHRRepo(s.db)
	}
	return s.HRS
}

func (s *Storage) BlackList() storage.BlackListI {
	if s.BlackListS == nil {
		s.BlackListS = NewBalckListRepo(s.db)
	}
	return s.BlackListS
}
