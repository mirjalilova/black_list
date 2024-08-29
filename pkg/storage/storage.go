package storage

import (
	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type StorageI interface {
	Admin() AdminI
	HR() HRI
	BlackList() BlackListI
}

type AdminI interface {
	Approve(id *pb.CreateHR) (*pb.Void, error)
	ListHR(filter *pb.Filter) (*pb.GetAllHRRes, error)
	Delete(id *pb.GetById) (*pb.Void, error)
	GetAllUsers(req *pb.ListUserReq) (*pb.ListUserRes, error)
	ChangeRole(req *pb.ChangeRoleReq) (*pb.Void, error)
}

type HRI interface {
	Create(req *pb.EmployeeCreate) (*pb.Void, error)
	Get(id *pb.GetById) (*pb.Employee, error)
	GetAll(req *pb.ListEmployeeReq) (*pb.ListEmployeeRes, error)
	Update(req *pb.UpdateReq) (*pb.Void, error)
	Delete(id *pb.GetById) (*pb.Void, error)
}

type BlackListI interface {
	Add(req *pb.BlackListCreate) (*pb.Void, error)
	Remove(id *pb.RemoveReq) (*pb.Void, error)
	GetAll(req *pb.Filter) (*pb.GetAllBlackListRes, error)
	MonitoringDailyReport(req *pb.Filter) (*pb.Reports, error)
	MonitoringWeeklyReport(req *pb.Filter) (*pb.Reports, error)
	MonitoringMonthlyReport(req *pb.Filter) (*pb.Reports, error)
}