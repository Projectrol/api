package main

import (
	"context"
	"net/http"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	data := &pb.CreateUserRequest{}
	err := common.ReadJSON(r, data)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid data"})
	}
	response, err := app.UserGRPCClient.CreateUser(context.Background(), data)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
	}
	common.WriteJSON(w, http.StatusCreated, response)
}

func (app *application) CreateWSHandler(w http.ResponseWriter, r *http.Request) {
	data := &pb.CreateWorkspaceRequest{}
	err := common.ReadJSON(r, data)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": "invalid data"})
	}
	response, err := app.WSGRPCClient.CreateWorkspace(context.Background(), data)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, common.Envelop{"error": err.Error()})
	}
	common.WriteJSON(w, http.StatusCreated, response)
}
