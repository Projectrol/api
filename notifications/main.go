package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/notifications/models"
	_ "github.com/lib/pq"
	"github.com/resend/resend-go/v2"
	"google.golang.org/grpc"
)

func SendMail() {
	apiKey := "re_7zJ1XWch_Df9MSNkcK4Pn27GydncqhEX9"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"hoangvule100@gmail.com"},
		Subject: "Hello World",
		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Print(sent)
}

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
	db, err := sql.Open("postgres", fmt.Sprintf(`user=%s 
	password=%s dbname=%s host=%s port=%d binary_parameters=yes`,
		dbUsername,
		dbPassword,
		dbName,
		dbHost,
		port,
	))
	if err != nil {
		log.Fatalf("Cannot connect to database. Error: " + err.Error())
	}

	lis, err := net.Listen("tcp", "localhost:3002")
	if err != nil {
		log.Fatalf("Cannot listen to tcp localhost:3002. Error: " + err.Error())
	}
	grpcServer := grpc.NewServer()
	server := &server{
		UserNotificationsSettingsModel: models.NewUserNotificationsSettingsModel(db),
	}
	pb.RegisterNotificationsServiceServer(grpcServer, server)
	log.Print("Starting notifications gRPC server at localhost:3002")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("Cannot start notifications gRPC server. Error: " + err.Error())
	}
}
