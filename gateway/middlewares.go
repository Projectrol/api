package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (app *application) AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := r.Cookie("access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		tokenStr := cookieValue.Value
		claims, err := common.ParseToken(tokenStr, "access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		userId := -1
		for k, v := range claims {
			if k == "sub" {
				userId = int(v.(float64))
			}
		}
		if userId == -1 {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		ctx := context.WithValue(r.Context(), common.ContextUserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) AuthorizeGuard(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceTag := strings.Split((strings.Split(r.URL.Path, "/api/workspaces/")[1]), "/")[1]
		workspaceId, _ := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
		method := r.Method

		cookieValue, err := r.Cookie("access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		tokenStr := cookieValue.Value
		claims, err := common.ParseToken(tokenStr, "access_token")
		if err != nil {
			common.WriteJSON(w, http.StatusUnauthorized, common.Envelop{"error": "unauthorized user"})
			return
		}
		var workspacesRoleIdList []*pb.WorkspaceRoleId
		for k, v := range claims {
			if k == "workspaces_role" {
				a, _ := json.Marshal(v)
				_ = json.Unmarshal(a, &workspacesRoleIdList)
			}
		}

		conn, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			return
		}

		var roleId int32

		for _, workspaceRoleId := range workspacesRoleIdList {
			if workspaceRoleId.WorkspaceId == int32(workspaceId) {
				roleId = int32(workspaceRoleId.RoleId)
				break
			}
		}

		c := pb.NewWorkspacesServiceClient(conn)

		if resourceTag == "projects" {
			projectSlug := httprouter.ParamsFromContext(r.Context()).ByName("projectSlug")
			if projectSlug != "" {
				response, err := c.CheckRoleValidForResource(context.Background(), &pb.CheckRoleValidForResourceRequest{
					RoleId:      int32(roleId),
					ResourceTag: resourceTag,
					Method:      method,
				})
				if err != nil || !response.IsValid {
					common.WriteJSON(w, http.StatusForbidden, common.Envelop{"error": "user don't have permission to access this resource"})
					return
				}
				next.ServeHTTP(w, r)
				return
			}
		}

		response, err := c.CheckRoleValidForResource(context.Background(), &pb.CheckRoleValidForResourceRequest{
			RoleId:      int32(roleId),
			ResourceTag: resourceTag,
			Method:      method,
		})

		if err != nil || !response.IsValid {
			common.WriteJSON(w, http.StatusForbidden, common.Envelop{"error": "user don't have permission to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
