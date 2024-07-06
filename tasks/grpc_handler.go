package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/tasks/models"
)

type server struct {
	pb.UnimplementedTasksServiceServer
	TaskModel *models.TaskModel
}

func (s *server) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	return s.TaskModel.Insert(ctx, in)
}
