package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
)

type HRService struct {
	storage storage.StorageI
	pb.UnimplementedHRServiceServer
}

func NewHRService(storage storage.StorageI) *HRService {
	return &HRService{
		storage: storage,
	}
}

func (s *HRService) Create(c context.Context, req *pb.EmployeeCreate) (*pb.Void, error) {
	return s.storage.HR().Create(req)
}

func (s *HRService) Get(c context.Context, id *pb.GetById) (*pb.Employee, error) {
	return s.storage.HR().Get(id)
}

func (s *HRService) GetAll(c context.Context, req *pb.ListEmployeeReq) (*pb.ListEmployeeRes, error) {
	return s.storage.HR().GetAll(req)
}

func (s *HRService) Update(c context.Context, req *pb.UpdateReq) (*pb.Void, error) {
	return s.storage.HR().Update(req)
}

func (s *HRService) Delete(c context.Context, id *pb.GetById) (*pb.Void, error) {
	return s.storage.HR().Delete(id)
}
