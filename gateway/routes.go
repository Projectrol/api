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

	router.HandlerFunc(http.MethodPost, "/api/calendar-events/create", app.AuthGuard(app.CreateCalendarEventHandler))

	router.HandlerFunc(http.MethodPost, "/api/projects/create", app.AuthGuard(app.CreateProjectHandler))
	router.HandlerFunc(http.MethodGet, "/api/projects/:workspaceSlug/:projectSlug", app.AuthGuard(app.GetProjectDetailsHandler))

	router.HandlerFunc(http.MethodGet, "/api/notifications/user-settings", app.AuthGuard(app.GetUserNotificationsSettingsHandler))
	router.HandlerFunc(http.MethodPatch, "/api/notifications/user-settings", app.AuthGuard(app.UpdateUserNotificationsSettingsHandler))

	handler := crs.Handler(router)
	return &handler
}
