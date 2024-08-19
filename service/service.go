package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
)

type ServiceI interface {
	Admin() AdminI
	HR() HRI
	BlackList() BlackListI
}

type AdminI interface {
	Approve(ctx context.Context, request *pb.CreateHR) (*pb.Void, error)
	ListHR(ctx context.Context, request *pb.Filter) (*pb.GetAllHRRes, error)
	Delete(ctx context.Context, request *pb.GetById) (*pb.Void, error)
}

type HRI interface {
	Create(ctx context.Context, request *pb.EmployeeCreate) (*pb.Void, error)
	Update(ctx context.Context, request *pb.UpdateReq) (*pb.Void, error)
	Delete(ctx context.Context, request *pb.GetById) (*pb.ListEmployeeRes, error)
	Get(ctx context.Context, request *pb.GetById) (*pb.Employee, error)
	GetAll(ctx context.Context, request *pb.ListEmployeeReq) (*pb.ListEmployeeReq, error)
}

type BlackListI interface {
	Add(ctx context.Context, request *pb.BlackListCreate) (*pb.Void, error)
	GetAll(ctx context.Context, request *pb.Filter) (*pb.Void, error)
	Remove(ctx context.Context, request *pb.RemoveReq) (*pb.Void, error)
}