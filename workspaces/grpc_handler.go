package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	m "github.com/lehoangvuvt/projectrol/users/models"
)

type server struct {
	pb.UnimplementedWorkspacesServiceServer
	models *m.Models
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {

	return nil, nil
}
