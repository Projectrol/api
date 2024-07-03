package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/users/models"
)

type server struct {
	pb.UnimplementedUsersServiceServer
	UserModel *models.UserModel
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return s.UserModel.Insert(ctx, in)
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.User, error) {
	return s.UserModel.Login(ctx, in)
}

func (s *server) GetUserById(ctx context.Context, in *pb.GetUserByIdRequest) (*pb.User, error) {
	return s.UserModel.GetUserById(ctx, in)
}
