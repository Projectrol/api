package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/workspaces/models"
)

type server struct {
	pb.UnimplementedWorkspacesServiceServer
	WorkspaceModel     *models.WorkspaceModel
	CalendarEventModel *models.CalendarEventModel
	ProjectModel       *models.ProjectsModel
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

func (s *server) CreateTask(ctx context.Context, in *pb.CreateCalendarEventRequest) (*pb.CreateCalendarEventResponse, error) {
	return s.CalendarEventModel.Insert(ctx, in)
}

func (s *server) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	return s.ProjectModel.CreateProject(ctx, in)
}

func (s *server) GetProjectsByWorkspaceId(ctx context.Context, in *pb.GetProjectsByWorkspaceIdRequest) (*pb.GetProjectsByWorkspaceIdResponse, error) {
	return s.ProjectModel.GetProjectsByWorkspaceId(ctx, in)
}

func (s *server) GetProjectDetails(ctx context.Context, in *pb.GetProjectDetailsRequest) (*pb.GetProjectDetailsResponse, error) {
	return s.ProjectModel.GetProjectDetails(ctx, in)
}

func (s *server) GetWorkspaceDetails(ctx context.Context, in *pb.GetWorkspaceDetailsRequest) (*pb.GetWorkspaceDetailsResponse, error) {
	return s.WorkspaceModel.GetWorkspaceDetails(ctx, in)
}

func (s *server) GetPermissions(ctx context.Context, in *pb.EmptyRequest) (*pb.GetPermissionsResponse, error) {
	return s.WorkspaceModel.GetPermissions(ctx, in)
}

func (s *server) UpdateRolePermission(ctx context.Context, in *pb.UpdateRolePermissionRequest) (*pb.UpdateRolePermissionResponse, error) {
	return s.WorkspaceModel.UpdateRolePermission(ctx, in)
}

func (s *server) GetWorkspaceRoles(ctx context.Context, in *pb.GetWorkspaceRolesRequest) (*pb.GetWorkspaceRolesResponse, error) {
	return s.WorkspaceModel.GetWorkspaceRoles(ctx, in)
}

func (s *server) GetUserRoleInWorkspace(ctx context.Context, in *pb.GetUserRoleInWorkspaceRequest) (*pb.GetUserRoleInWorkspaceResponse, error) {
	return nil, nil
}
