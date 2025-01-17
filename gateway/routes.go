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

	router.HandlerFunc(http.MethodPost, "/api/workspaces", app.AuthGuard(app.CreateWorkspaceHandler))

	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/search-projects/:q", app.AuthGuard(app.AuthorizeGuard(app.GetProjectsByWorkspaceIdHandler)))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id", app.AuthGuard(app.GetWorkspaceDetailsHandler))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/roles", app.AuthGuard(app.AuthorizeGuard(app.GetWorkspaceRolesHandler)))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/user/role", app.AuthGuard(app.GetUserRoleInWorkspaceHandler))

	router.HandlerFunc(http.MethodPost, "/api/calendar-events/create", app.AuthGuard(app.CreateCalendarEventHandler))

	router.HandlerFunc(http.MethodPost, "/api/workspaces/:id/projects/create", app.AuthGuard(app.AuthorizeGuard(app.CreateProjectHandler)))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/projects/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.GetProjectDetailsHandler)))
	router.HandlerFunc(http.MethodDelete, "/api/workspaces/:id/projects/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.DeleteProjectHandler)))

	router.HandlerFunc(http.MethodGet, "/api/notifications/user-settings", app.AuthGuard(app.GetUserNotificationsSettingsHandler))
	router.HandlerFunc(http.MethodPatch, "/api/notifications/user-settings", app.AuthGuard(app.UpdateUserNotificationsSettingsHandler))

	router.HandlerFunc(http.MethodGet, "/api/permissions", app.GetPermissionsHandler)

	router.HandlerFunc(http.MethodPost, "/api/workspaces/:id/roles", app.AuthGuard(app.AuthorizeGuard(app.CreateNewRoleHandler)))
	router.HandlerFunc(http.MethodPatch, "/api/workspaces/:id/roles", app.AuthGuard(app.AuthorizeGuard(app.UpdateRolePermissionHandler)))

	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/members", app.AuthGuard(app.AuthorizeGuard(app.GetWorkspaceMembersHandler)))

	router.HandlerFunc(http.MethodPost, "/api/workspaces/:id/tasks/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.CreateTaskHandler)))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/tasks/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.GetProjectTasksHandler)))

	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/tasks/:projectSlug/:nanoid", app.AuthGuard(app.AuthorizeGuard(app.GetProjectTaskDetailsHandler)))
	router.HandlerFunc(http.MethodPatch, "/api/workspaces/:id/tasks/:projectSlug/:nanoid", app.AuthGuard(app.AuthorizeGuard(app.UpdateProjectTaskHandler)))

	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/documents/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.GetProjectDocumentsHandler)))
	router.HandlerFunc(http.MethodPost, "/api/workspaces/:id/documents/:projectSlug", app.AuthGuard(app.AuthorizeGuard(app.CreateProjectDocumentHandler)))
	router.HandlerFunc(http.MethodGet, "/api/workspaces/:id/documents/:projectSlug/:nanoid", app.AuthGuard(app.AuthorizeGuard(app.GetProjectDocumentDetailsHandler)))
	router.HandlerFunc(http.MethodPatch, "/api/workspaces/:id/documents/:projectSlug/:nanoid", app.AuthGuard(app.AuthorizeGuard(app.UpdateProjectDocumentDetailsHandler)))

	handler := crs.Handler(router)
	return &handler
}
