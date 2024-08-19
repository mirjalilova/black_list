package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
)

type BlackListService struct {
	storage storage.StorageI
	pb.UnimplementedBlackListServiceServer
}

func NewBlackListService(storage storage.StorageI) *BlackListService {
	return &BlackListService{
		storage: storage,
	}
}

func (s *BlackListService) Add(c context.Context, req *pb.BlackListCreate) (*pb.Void, error) {
	return s.storage.BlackList().Add(req)
}

func (s *BlackListService) GetAll(c context.Context, req *pb.Filter) (*pb.GetAllBlackListRes, error) {
	return s.storage.BlackList().GetAll(req)
}

func (s *BlackListService) Remove(c context.Context, req *pb.RemoveReq) (*pb.Void, error) {
	return s.storage.BlackList().Remove(req)
}
