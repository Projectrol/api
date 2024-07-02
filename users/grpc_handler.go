package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	m "github.com/lehoangvuvt/projectrol/users/models"
)

type server struct {
	pb.UnimplementedUsersServiceServer
	models *m.Models
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	err := s.models.UserModel.Insert(in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		Success: true,
		Message: "create user success",
	}, nil
}
