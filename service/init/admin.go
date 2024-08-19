package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
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
