package main

import (
	"context"
	"log"
	"net/http"

	common "github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to users gRPC server. Error: " + err.Error())
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
		log.Fatal("Cannot connect to users gRPC server. Error: " + err.Error())
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
		log.Fatal("Cannot connect to users gRPC server. Error: " + err.Error())
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
		log.Fatal("Cannot connect to workspaces gRPC server. Error: " + err.Error())
	}
	client := pb.NewWorkspacesServiceClient(conn)
	request := &pb.CreateWorkspaceRequest{}
	err = common.ReadJSON(r, request)

	if err != nil {
		errMsg := common.Envelop{"error": "error input format"}
		common.WriteJSON(w, http.StatusBadRequest, errMsg)
		return
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
