package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) getRoutes() *http.Handler {
	router := httprouter.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "PATCH"},
	})

	router.HandlerFunc(http.MethodPost, "/api/users/create", app.CreateUserHandler)
	router.HandlerFunc(http.MethodPatch, "/api/users/settings", app.AuthGuard(app.UpdateUserSettingsHandler))
	router.HandlerFunc(http.MethodGet, "/api/users/workspaces", app.AuthGuard(app.GetWorkspacesByUserId))
	router.HandlerFunc(http.MethodPost, "/api/users/login", app.LoginHandler)
	router.HandlerFunc(http.MethodGet, "/api/users/logout", app.AuthGuard(app.Logout))
	router.HandlerFunc(http.MethodGet, "/api/authenticate", app.AuthGuard(app.AuthenticateHandler))

	router.HandlerFunc(http.MethodPost, "/api/workspaces/create", app.AuthGuard(app.CreateWorkspaceHandler))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/projects", app.AuthGuard(app.GetProjectsByWorkspaceIdHandler))

	router.HandlerFunc(http.MethodPost, "/api/tasks/create", app.AuthGuard(app.CreateTaskHandler))

	router.HandlerFunc(http.MethodPost, "/api/projects/create", app.AuthGuard(app.CreateProjectHandler))

	handler := crs.Handler(router)
	return &handler
}
