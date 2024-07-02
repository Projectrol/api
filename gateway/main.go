package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type application struct {
	UserGRPCClient pb.UsersServiceClient
	WSGRPCClient   pb.WorkspacesServiceClient
}

func main() {
	app := &application{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn2, err := grpc.NewClient("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial user gRPC server: %v", err)
	}
	defer conn.Close()

	userServerClient := pb.NewUsersServiceClient(conn)
	wsServerClient := pb.NewWorkspacesServiceClient(conn2)

	app.UserGRPCClient = userServerClient
	app.WSGRPCClient = wsServerClient

	server := &http.Server{
		Addr:     ":8080",
		Handler:  app.getRoutes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	log.Printf("Listening on port %d", 8080)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Cannot start gateway service. Error: " + err.Error())
		os.Exit(1)
	}
}
