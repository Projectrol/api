package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	pm "github.com/lehoangvuvt/projectrol/projects/models"
)

type server struct {
	pb.UnimplementedProjectsServiceServer
	ProjectsModel *pm.ProjectsModel
}

func (s *server) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	return s.ProjectsModel.CreateProject(ctx, in)
}

func (s *server) GetProjectsByWorkspaceId(ctx context.Context, in *pb.GetProjectsByWorkspaceIdRequest) (*pb.GetProjectsByWorkspaceIdResponse, error) {
	return s.ProjectsModel.GetProjectsByWorkspaceId(ctx, in)
}
