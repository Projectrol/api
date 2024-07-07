package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	common "github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
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
	token, _ := common.SignToken(common.Envelop{"sub": res.Id}, 60, "access_token")
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

func (app *application) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg := common.Envelop{"error": "Cannot connect to workspaces gRPC server. Error: " + err.Error()}
		common.WriteJSON(w, http.StatusInternalServerError, errMsg)
		return
	}
	client := pb.NewWorkspacesServiceClient(conn)
	body := &wsModels.CreateTaskInput{}
	userId := r.Context().Value(common.ContextUserIdKey).(int)
	err = common.ReadJSON(r, body)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	request := &pb.CreateTaskRequest{
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
	res, err := client.CreateTask(ctx, request)
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

	ctx := context.Background()
	res, err := client.GetProjectsByWorkspaceId(ctx, request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusCreated, common.Envelop{"projects": res.Projects})
}

func (app *application) GetProjectDetails(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	workspaceSlug := params.ByName("workspaceSlug")
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
		WorkspaceSlug: workspaceSlug,
		ProjectSlug:   projectSlug,
		UserId:        int32(userId),
	}
	response, err := client.GetProjectDetails(context.Background(), request)
	if err != nil {
		errMsg := common.Envelop{"error": err.Error()}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	common.WriteJSON(w, http.StatusOK, common.Envelop{"project": response.Project})
}
