package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
	"golang.org/x/exp/slog"
)

type AdminService struct {
	storage storage.StorageI
	pb.UnimplementedAdminServiceServer
}

func NewAdminService(storage storage.StorageI) *AdminService {
	return &AdminService{
		storage: storage,
	}
}

func (s *AdminService) Approve(c context.Context, id *pb.CreateHR) (*pb.Void, error) {
	return s.storage.Admin().Approve(id)
}

func (s *AdminService) ListHR(c context.Context, filter *pb.Filter) (*pb.GetAllHRRes, error) {
    return s.storage.Admin().ListHR(filter)
}

func (s *AdminService) Delete(c context.Context, req *pb.GetById) (*pb.Void, error) {
    return s.storage.Admin().Delete(req)
}

func (s *AdminService) GetAllUsers(c context.Context, req *pb.ListUserReq) (*pb.ListUserRes, error) {
    res, err := s.storage.Admin().GetAllUsers(req)
    if err!= nil {
        slog.Error("Error getting all users: %v", err)
        return nil, err
    }

    slog.Info("Got all users: %+v", res)
    return res, nil
}

func (s *AdminService) ChangeRole(c context.Context, req *pb.ChangeRoleReq) (*pb.Void, error) {
    return s.storage.Admin().ChangeRole(req)
}