package service

import (
	"context"

	pb "github.com/mirjalilova/black_list/internal/genproto/black_list"
	"github.com/mirjalilova/black_list/pkg/storage"
	"golang.org/x/exp/slog"
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
	_, err := s.storage.BlackList().Add(req)
	if err!= nil {
		slog.Error("Error adding black list: %v", err)
        return nil, err
    }

	slog.Info("Black list added")
	return &pb.Void{}, nil
}

func (s *BlackListService) GetAll(c context.Context, req *pb.Filter) (*pb.GetAllBlackListRes, error) {
	res, err := s.storage.BlackList().GetAll(req)
	if err!= nil {
        slog.Error("Error getting black list: %v", err)
        return nil, err
    }

	slog.Info("Got black list: %+v", res)
	return res, nil
}

func (s *BlackListService) Remove(c context.Context, req *pb.RemoveReq) (*pb.Void, error) {
	_, err := s.storage.BlackList().Remove(req)
	if err!= nil {
        slog.Error("Error removing black list: %v", err)
        return nil, err
    }

	slog.Info("Black list removed")
	return &pb.Void{}, nil
}
