package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/users/models"
	"google.golang.org/grpc"
)

var (
	GRPC_PORT = 3000
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", GRPC_PORT))
	if err != nil {
		log.Fatal("Cannot start gRPC server. Error: " + err.Error())
	}
	grpcServer := grpc.NewServer()
	db, err := common.ConnectToDB()
	if err != nil {
		log.Fatal("Cannot connect to database. Error: " + err.Error())
	}
	server := &server{models: models.NewModels(db)}
	pb.RegisterUsersServiceServer(grpcServer, server)
	log.Printf("Starting users gRPC server at: %d", lis.Addr())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
