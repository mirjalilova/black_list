package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
	"golang.org/x/exp/slog"
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
	_, err := s.storage.HR().Create(req)
	if err!= nil {
		slog.Error("Error creating Employee: %v", err)
        return nil, err
    }

	slog.Info("Employee created")
	return &pb.Void{}, nil
}

func (s *HRService) Get(c context.Context, id *pb.GetById) (*pb.Employee, error) {
	res, err := s.storage.HR().Get(id)
	if err!= nil {
        slog.Error("Error getting Employee: %v", err)
        return nil, err
    }

	slog.Info("Got Employee: %+v", res)
	return res, nil
}

func (s *HRService) GetAll(c context.Context, req *pb.ListEmployeeReq) (*pb.ListEmployeeRes, error) {
	res, err := s.storage.HR().GetAll(req)
	if err!= nil {
        slog.Error("Error getting Employee: %v", err)
        return nil, err
    }

	slog.Info("Got Employee: %+v", res)
	return res, nil
}

func (s *HRService) Update(c context.Context, req *pb.UpdateReq) (*pb.Void, error) {
	_, err := s.storage.HR().Update(req)
	if err!= nil {
		slog.Error("Error updating Employee: %v", err)
        return nil, err
    }

    slog.Info("Employee updated")
	return &pb.Void{}, nil
}

func (s *HRService) Delete(c context.Context, id *pb.GetById) (*pb.Void, error) {
	_, err := s.storage.HR().Delete(id)
	if err!= nil {
        slog.Error("Error deleting Employee: %v", err)
        return nil, err
    }

	slog.Info("Employee deleted")
	return &pb.Void{}, nil
}