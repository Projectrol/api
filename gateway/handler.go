package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	common "github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	notiModels "github.com/lehoangvuvt/projectrol/notifications/models"
	wsModels "github.com/lehoangvuvt/projectrol/workspaces/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to users gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewUsersServiceClient(conn)
	request := &pb.CreateUserRequest{}
	err = common.ReadJSON(r, request)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	ctx := context.Background()
	res, err := client.CreateUser(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	successMsg := common.Envelop{"message": "success", "id": res.Id}
	common.WriteJSON(w, http.StatusCreated, successMsg)
}

func (app *application) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to users gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewUsersServiceClient(conn)
	request := &pb.GetUserByIdRequest{UserId: int32(userId)}
	user, err := client.GetUserById(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": "unauthorized user"}
		common.WriteJSON(w, http.StatusUnauthorized, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"user": user})
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to users gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewUsersServiceClient(conn)
	request := &pb.LoginRequest{}
	err = common.ReadJSON(r, request)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	ctx := context.Background()
	res, err := client.Login(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusUnauthorized, errMsg)
		return
	}
	successMsg := common.Envelop{"message": "success", "user": res}
	dataToSign := common.Envelop{"sub": res.Id, "workspaces_role": res.WorkspaceRoleIdList}
	token, _ := common.SignToken(dataToSign, 60, "access_token")

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		MaxAge:   60 * 60,
	})
	common.WriteJSON(w, http.StatusCreated, successMsg)
}

func (app *application) CreateWorkspaceHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	body := &wsModels.CreateWorkspaceInput{}
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	err = common.ReadJSON(r, body)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	request := &pb.CreateWorkspaceRequest{
		Name:    body.Name,
		OwnerId: int32(userId),
		Logo:    body.Logo,
	}

	ctx := context.Background()
	res, err := client.CreateWorkspace(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	successMsg := common.Envelop{"message": "success", "id": res.Nanoid}
	common.WriteJSON(w, http.StatusCreated, successMsg)
}

func (app *application) GetWorkspacesByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, common.Envelop{"error": "cannot connect to workspace gRPC server"})
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.GetWorkspacesByUserId(context.Background(), &pb.GetWorkspacesByUserIdRequest{UserId: int32(userId)})
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, common.Envelop{"error": err.Error()})
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"workspaces": response.Workspaces})
}

func (app *application) CreateCalendarEventHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	body := &wsModels.CreateEventInput{}
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	err = common.ReadJSON(r, body)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	request := &pb.CreateCalendarEventRequest{
		Title:       body.Title,
		Description: body.Description,
		Duration:    int32(body.Duration),
		Dtstart:     body.DtStart,
		Type:        body.Type,
		Recurring:   body.Recurring,
		CreatedBy:   int32(userId),
		WorkspaceId: int32(body.WorkspaceId),
	}

	ctx := context.Background()
	res, err := client.CreateCalendarEvent(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	successMsg := common.Envelop{"message": "success", "id": res.Id}
	common.WriteJSON(w, http.StatusCreated, successMsg)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		MaxAge:   0,
		Expires:  time.Now(),
	})
}

func (app *application) UpdateUserSettingsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to users gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewUsersServiceClient(conn)
	request := &pb.UserSettings{}
	err = common.ReadJSON(r, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	request.Id = int32(userId)
	userSettings, err := client.UpdateUserSettings(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"settings": userSettings})
}

func (app *application) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	body := &wsModels.CreateProjectInput{}
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	err = common.ReadJSON(r, body)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	request := &pb.CreateProjectRequest{
		WorkspaceId: int32(body.WorkspaceId),
		CreatedBy:   int32(userId),
		Name:        body.Name,
		Summary:     body.Summary,
		Description: body.Description,
		Dtstart:     body.Dtstart,
		Dtend:       body.Dtend,
		IsPrivate:   body.IsPrivate,
	}
	ctx := context.Background()
	res, err := client.CreateProject(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	successMsg := common.Envelop{"message": "success", "id": res.Id}
	common.WriteJSON(w, http.StatusCreated, successMsg)
}

func (app *application) GetProjectsByWorkspaceIdHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	workspaceId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || workspaceId < 0 {
		errMsg := common.Envelop{"error": "Invalid workspace id"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)

	request := &pb.GetProjectsByWorkspaceIdRequest{
		WorkspaceId: int32(workspaceId),
		UserId:      int32(userId),
	}

	q := params.ByName("q")
	if q != "" {
		request.Q = q
	} else {
		request.Q = "*"
	}

	ctx := context.Background()
	res, err := client.GetProjectsByWorkspaceId(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"projects": res.Projects})
}

func (app *application) GetWorkspaceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	workspaceId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || workspaceId < 0 {
		errMsg := common.Envelop{"error": "Invalid workspace id"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	// userId := r.Context().Value(common.ContextUserIdKey).(int)
	request := &pb.GetWorkspaceDetailsRequest{
		Id: int32(workspaceId),
	}

	ctx := context.Background()
	response, err := client.GetWorkspaceDetails(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response})
}

func (app *application) GetProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	projectSlug := params.ByName("projectSlug")
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	request := &pb.GetProjectDetailsRequest{
		ProjectSlug: projectSlug,
		UserId:      int32(userId),
	}
	response, err := client.GetProjectDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response.Details})
}

func (app *application) GetUserNotificationsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to notifications gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewNotificationsServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	request := &pb.GetUserNotificationsSettingsRequest{
		UserId: int32(userId),
	}
	response, err := client.GetUserNotificationsSettings(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"settings": response.Settings})
}

func (app *application) UpdateUserNotificationsSettingsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to notifications gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewNotificationsServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	bodyData := &notiModels.UpsertUserNotiSettingsInput{}
	err = common.ReadJSON(r, bodyData)
	if err != nil {
		errMsg := common.Envelop{"error": "Invalid body request data. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	request := &pb.UpsertUserNotificationsSettingsRequest{
		Settings: &pb.UserNotificationsSettings{
			UserId:              int32(userId),
			IsViaInbox:          bodyData.IsViaInbox,
			IsViaEmail:          bodyData.IsViaEmail,
			TaskNotiSettings:    bodyData.TaskNotiSettings,
			ProjectNotiSettings: bodyData.ProjectNotiSettings,
			EventNotiSettings:   bodyData.EventNotiSettings,
			EventNoticeBefore:   int32(bodyData.EventNoticeBefore),
		},
	}
	response, err := client.UpdateUserNotificationsSettings(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"settings": response.Settings})
}

func (app *application) GetPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	request := &pb.EmptyRequest{}

	ctx := context.Background()
	response, err := client.GetPermissions(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusCreated, common.Envelop{"permissions": response.Permissions})
}

func (app *application) CreateNewRoleHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	c := pb.NewWorkspacesServiceClient(conn)
	request := &pb.CreateNewRoleRequest{}
	err = common.ReadJSON(r, request)
	if err != nil {
		errMsg := common.Envelop{"error": "Invalid body data. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	response, err := c.CreateNewRole(context.Background(), request)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
		return
	}
	common.WriteJSON(w, http.StatusCreated, common.Envelop{"role_id": response.Id})
}

func (app *application) UpdateRolePermissionHandler(w http.ResponseWriter, r *http.Request) {
	request := &pb.UpdateRolePermissionRequest{}
	err := common.ReadJSON(r, request)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid body format"})
		return
	}
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	c := pb.NewWorkspacesServiceClient(conn)
	response, err := c.UpdateRolePermission(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"role": response.Role})
}

func (app *application) GetWorkspaceRolesHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	workspaceId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || workspaceId < 0 {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid workspace id"})
		return
	}
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	request := &pb.GetWorkspaceRolesRequest{
		Id: int32(workspaceId),
	}
	c := pb.NewWorkspacesServiceClient(conn)
	response, err := c.GetWorkspaceRoles(context.Background(), request)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"roles": response.Roles})
}

func (app *application) GetUserRoleInWorkspaceHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	workspaceId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || workspaceId < 0 {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid workspace id"})
		return
	}
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	c := pb.NewWorkspacesServiceClient(conn)
	request := &pb.GetUserRoleInWorkspaceRequest{
		UserId:      int32(userId),
		WorkspaceId: int32(workspaceId),
	}
	response, err := c.GetUserRoleInWorkspace(context.Background(), request)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"role": response.Role})
}

func (app *application) GetWorkspaceMembersHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	worksapceId, err := strconv.Atoi(params.ByName("id"))
	if err != nil || worksapceId < 0 {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid workspace id"})
		return
	}
	request := &pb.GetWorkspaceMembersRequest{
		Id: int32(worksapceId),
	}
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	c := pb.NewWorkspacesServiceClient(conn)
	response, err := c.GetWorkspaceMembers(context.Background(), request)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"members": response.Members})
}

func (app *application) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	projectSlug := params.ByName("projectSlug")
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	request := &pb.GetProjectDetailsRequest{
		ProjectSlug: projectSlug,
		UserId:      int32(userId),
	}
	response, err := client.GetProjectDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response.Details})
}

func (app *application) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)

	input := &wsModels.CreateTaskInput{}
	err = common.ReadJSON(r, input)

	if err != nil {
		log.Print(err)
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid body format"})
		return
	}

	request := &pb.CreateTaskRequest{
		UserId:      int32(userId),
		ProjectSlug: input.ProjectSlug,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Label:       input.Label,
		Priority:    input.Priority,
		IsPublished: input.IsPublished,
	}

	response, err := client.CreateTask(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusCreated, common.Envelop{"nanoid": response.Nanoid})
}

func (app *application) GetProjectTasksHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	projectSlug := httprouter.ParamsFromContext(r.Context()).ByName("projectSlug")
	client := pb.NewWorkspacesServiceClient(conn)
	userId := r.Context().Value(common.ContextUserIdKey).(int)

	request := &pb.GetProjectTasksRequest{
		UserId:      int32(userId),
		ProjectSlug: projectSlug,
	}

	response, err := client.GetProjectTasks(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"tasks": response.Tasks})
}

func (app *application) GetProjectDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	projectSlug := httprouter.ParamsFromContext(r.Context()).ByName("projectSlug")

	request := &pb.GetProjectDocumentsRequest{
		ProjectSlug: projectSlug,
	}

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.GetProjectDocuments(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"documents": response.Documents})
}

func (app *application) CreateProjectDocumentHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	projectSlug := httprouter.ParamsFromContext(r.Context()).ByName("projectSlug")

	input := &wsModels.CreateProjectDocumentInput{}
	err := common.ReadJSON(r, input)

	if err != nil {
		log.Print(err)
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid body format"})
		return
	}

	request := &pb.CreateProjectDocumentRequest{
		UserId:      int32(userId),
		ProjectSlug: projectSlug,
		Name:        input.Name,
		Content:     input.Content,
	}

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.CreateProjectDocument(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusCreated, common.Envelop{"nanoid": response.Nanoid})
}

func (app *application) GetProjectDocumentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	nanoid := httprouter.ParamsFromContext(r.Context()).ByName("nanoid")

	request := &pb.GetProjectDocumentDetailsRequest{
		Nanoid: nanoid,
	}

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.GetProjectDocumentDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response.Details})
}

func (app *application) UpdateProjectDocumentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	nanoid := httprouter.ParamsFromContext(r.Context()).ByName("nanoid")
	userId := r.Context().Value(common.ContextUserIdKey).(int)

	input := &wsModels.UpdateProjectDocumentInput{}
	err := common.ReadJSON(r, input)

	if err != nil {
		log.Print(err)
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid body format"})
		return
	}

	request := &pb.UpdateProjectDocumentDetailsRequest{
		Nanoid:    nanoid,
		UpdatedBy: int32(userId),
		Name:      input.Name,
		Content:   input.Content,
	}

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.UpdateProjectDocumentDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response.Details})
}

func (app *application) UpdateProjectTaskHandler(w http.ResponseWriter, r *http.Request) {
	nanoid := httprouter.ParamsFromContext(r.Context()).ByName("nanoid")
	userId := r.Context().Value(common.ContextUserIdKey).(int)

	input := &wsModels.UpdateTaskInput{}
	err := common.ReadJSON(r, input)

	if err != nil {
		log.Print(err)
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid body format"})
		return
	}

	request := &pb.UpdateTaskRequest{
		Nanoid:       nanoid,
		UserId:       int32(userId),
		ChangedField: input.ChangedField,
		Value:        input.Value,
	}

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.UpdateTask(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"task": response.Task})
}

func (app *application) GetProjectTaskDetailsHandler(w http.ResponseWriter, r *http.Request) {
	nanoid := httprouter.ParamsFromContext(r.Context()).ByName("nanoid")

	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}

	request := &pb.GetProjectTaskDetailsRequest{
		Nanoid: nanoid,
	}

	client := pb.NewWorkspacesServiceClient(conn)
	response, err := client.GetProjectTaskDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	common.WriteJSON(w, http.StatusOK, common.Envelop{"details": response})
}
