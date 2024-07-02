package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) getRoutes() http.Handler {
	router := httprouter.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	router.HandlerFunc(http.MethodPost, "/api/users/create", app.CreateUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/workspaces/create", app.CreateWSHandler)

	handler := crs.Handler(router)

	return handler
}
