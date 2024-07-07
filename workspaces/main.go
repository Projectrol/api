package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/workspaces/models"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dbUsername := common.GetEnvValByKey("DB_USERNAME", "postgres")
	dbPassword := common.GetEnvValByKey("DB_PASSWORD", "admin")
	dbHost := common.GetEnvValByKey("DB_HOST", "localhost")
	dbName := common.GetEnvValByKey("DB_NAME", "postgres")
	portStr := common.GetEnvValByKey("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432
	}
	connectionStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dbUsername, dbPassword, dbHost, port, dbName)
	db, err := common.ConnectToDatabase(context.Background(), "postgres", connectionStr)
	if err != nil {
		log.Fatalf("Cannot connect to database. Error: " + err.Error())
	}
	lis, err := net.Listen("tcp", "localhost:3001")
	if err != nil {
		log.Fatalf("Cannot listen TCP on address %s", "localhost:3001")
	}
	grpcServer := grpc.NewServer()
	server := &server{
		WorkspaceModel: models.NewWorkspaceModel(db),
		TaskModel:      models.NewTaskModel(db),
		ProjectModel:   models.NewProjectsModel(db),
	}
	pb.RegisterWorkspacesServiceServer(grpcServer, server)
	log.Printf("Start gRPC server on address %s", "localhost:3001")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Cannot server gGRPC on port address %s", "localhost:3001")
	}
}
