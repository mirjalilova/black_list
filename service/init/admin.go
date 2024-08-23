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
	_, err := s.storage.Admin().Approve(id)
	if err!= nil {
        slog.Error("Error approving HR: %v", err)
        return nil, err
    }

	slog.Info("HR approved")
	return &pb.Void{}, nil
}

func (s *AdminService) ListHR(c context.Context, filter *pb.Filter) (*pb.GetAllHRRes, error) {
    res, err := s.storage.Admin().ListHR(filter)
	if err!= nil {
        slog.Error("Error getting HR: %v", err)
        return nil, err
    }

	slog.Info("Got HR: %+v", res)
	return res, nil
}

func (s *AdminService) Delete(c context.Context, req *pb.GetById) (*pb.Void, error) {
    _, err := s.storage.Admin().Delete(req)
	if err!= nil {
        slog.Error("Error deleting HR: %v", err)
        return nil, err
    }

	slog.Info("HR deleted")
	return &pb.Void{}, nil
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
    _, err := s.storage.Admin().ChangeRole(req)
    if err!= nil {
        slog.Error("Error changing role: %v", err)
        return nil, err
    }

    slog.Info("Role changed")
    return &pb.Void{}, nil
}

func (s *AdminService) MonitoringDailyReport(c context.Context, req *pb.Void) (*pb.Reports, error) {
	_, err := s.storage.BlackList().MonitoringDailyReport(req)
	if err!= nil {
        slog.Error("Error getting daily report: %v", err)
        return nil, err
    }

	slog.Info("Got daily report")
	return &pb.Reports{}, nil
}


func (s *AdminService) MonitoringWeeklyReport(c context.Context, req *pb.Void) (*pb.Reports, error) {
    _, err := s.storage.BlackList().MonitoringWeeklyReport(req)
	if err!= nil {
        slog.Error("Error getting weekly report: %v", err)
        return nil, err
    }

	slog.Info("Got weekly report")
	return &pb.Reports{}, nil
}


func (s *AdminService) MonitoringMonthlyReport(c context.Context, req *pb.Void) (*pb.Reports, error) {
    _, err := s.storage.BlackList().MonitoringMonthlyReport(req)
	if err!= nil {
        slog.Error("Error getting monthly report: %v", err)
        return nil, err
    }

	slog.Info("Got monthly report")
	return &pb.Reports{}, nil
}