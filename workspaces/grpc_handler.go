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
	TaskModel          *models.TaskModel
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

func (s *server) CreateCalendarEvent(ctx context.Context, in *pb.CreateCalendarEventRequest) (*pb.CreateCalendarEventResponse, error) {
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
	return s.WorkspaceModel.GetUserRoleInWorkspace(ctx, in)
}

func (s *server) GetRoleIdOfUserWorkspaces(ctx context.Context, in *pb.GetRoleIdOfUserWorkspacesRequest) (*pb.GetRoleIdOfUserWorkspacesResponse, error) {
	return s.WorkspaceModel.GetRoleIdOfUserWorkspaces(ctx, in)
}

func (s *server) CheckRoleValidForResource(ctx context.Context, in *pb.CheckRoleValidForResourceRequest) (*pb.CheckRoleValidForResourceResponse, error) {
	return s.WorkspaceModel.CheckRoleValidForResource(ctx, in)
}

func (s *server) CreateNewRole(ctx context.Context, in *pb.CreateNewRoleRequest) (*pb.CreateNewRoleResponse, error) {
	return s.WorkspaceModel.CreateNewRole(ctx, in)
}

func (s *server) GetWorkspaceMembers(ctx context.Context, in *pb.GetWorkspaceMembersRequest) (*pb.GetWorkspaceMembersResponse, error) {
	return s.WorkspaceModel.GetWorkspaceMembers(ctx, in)
}

func (s *server) CheckUserHasAccessToProject(ctx context.Context, in *pb.CheckUserHasAccessToProjectRequest) (*pb.CheckRoleValidForResourceResponse, error) {
	return s.WorkspaceModel.CheckUserHasAccessToProject(ctx, in)
}

func (s *server) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	return s.TaskModel.CreateTask(ctx, in)
}

func (s *server) GetProjectTasks(ctx context.Context, in *pb.GetProjectTasksRequest) (*pb.GetProjectTasksResponse, error) {
	return s.TaskModel.GetProjectTasks(ctx, in)
}

func (s *server) CreateProjectDocument(ctx context.Context, in *pb.CreateProjectDocumentRequest) (*pb.CreateProjectDocumentResponse, error) {
	return s.ProjectModel.CreateProjectDocument(ctx, in)
}

func (s *server) GetProjectDocuments(ctx context.Context, in *pb.GetProjectDocumentsRequest) (*pb.GetProjectDocumentsResponse, error) {
	return s.ProjectModel.GetProjectDocuments(ctx, in)
}

func (s *server) GetProjectDocumentDetails(ctx context.Context, in *pb.GetProjectDocumentDetailsRequest) (*pb.GetProjectDocumentDetailsResponse, error) {
	return s.ProjectModel.GetProjectDocumentDetails(ctx, in)
}

func (s *server) UpdateProjectDocumentDetails(ctx context.Context, in *pb.UpdateProjectDocumentDetailsRequest) (*pb.UpdateProjectDocumentDetailsResponse, error) {
	return s.ProjectModel.UpdateProjectDocumentDetails(ctx, in)
}
