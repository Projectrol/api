package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/workspaces/models"
)

type server struct {
	pb.UnimplementedWorkspacesServiceServer
	WorkspaceModel *models.WorkspaceModel
}

func (s *server) CreateWorkspace(ctx context.Context, in *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {
	nanoid, err := s.WorkspaceModel.Insert(ctx, in)
	if err != nil {
		return &pb.CreateWorkspaceResponse{
			Nanoid: "",
		}, err
	}
	return &pb.CreateWorkspaceResponse{
		Nanoid: nanoid,
	}, nil
}

func (s *server) GetWorkspacesByUserId(ctx context.Context, in *pb.GetWorkspacesByUserIdRequest) (*pb.GetWorkspacesByUserIdResponse, error) {
	return s.WorkspaceModel.GetWorkspacesByUserId(ctx, in.UserId)
}
