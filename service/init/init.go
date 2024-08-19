package service

import (
	st "github.com/mirjalilova/black_list/pkg/storage/postgres"
	"github.com/mirjalilova/black_list/service"
	"golang.org/x/exp/slog"
)

type Service struct {
	AdminS     service.AdminI
	HRS        service.HRI
	BlackListS service.BlackListI
}

func ConnectDB() (*st.Storage, error) {
	storage, err := st.ConnectDB()
	if err != nil {
		slog.Error("can't connect to postgres: %v", err)
		return nil, err
	}
	slog.Info("Connected to Postgres")
	return storage, nil
}
