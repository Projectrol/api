package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/tasks/models"
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
	server := &server{
		TaskModel: models.NewTaskModel(db),
	}
	lis, err := net.Listen("tcp", "localhost:3002")
	if err != nil {
		log.Fatal("Cannot listen to tcp server at localhost:3002. Error: " + err.Error())
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTasksServiceServer(grpcServer, server)
	log.Print("Starting tasks gRPC server at localhost:3002")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("Cannot server gRPC server. Error: " + err.Error())
	}
}
