package app

import (
	"net"

	"github.com/mirjalilova/black_list/internal/config"
	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	st "github.com/mirjalilova/black_list/service/init"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {
	conn, err := st.ConnectDB()
	if err != nil {
		slog.Error("Failed to connect to the database: %v", err)
		return
	}

	admin := st.NewAdminService(conn)
	black_list := st.NewBlackListService(conn)
	hr := st.NewHRService(conn)

	newServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(newServer, admin)
	pb.RegisterBlackListServiceServer(newServer, black_list)
	pb.RegisterHRServiceServer(newServer, hr)

	lis, err := net.Listen("tcp", cfg.BLACKLIST_PORT)
	if err != nil {
		slog.Error("Failed to listen on port %v: %v", cfg.BLACKLIST_PORT, err)
		return
	}
	slog.Info("Starting gRPC server on %v", cfg.BLACKLIST_PORT)
	if err := newServer.Serve(lis); err != nil {
		slog.Error("Failed to serve gRPC server: %v", err)
	}

}