package models

import (
	"context"
	"database/sql"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type TaskModel struct {
	DB *sql.DB
}

type CreateTaskInput struct {
	ProjectId   int    `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Label       string `json:"label"`
	IsPublished bool   `json:"is_published"`
}

func NewTaskModel(DB *sql.DB) *TaskModel {
	return &TaskModel{DB}
}

func (m *TaskModel) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	nanoid := common.GenerateNanoid(8)
	_, err := m.DB.Exec(`INSERT INTO tasks (nanoid, project_id, title, description, status, label, is_published, created_by) 
					VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		nanoid, in.ProjectId, in.Title, in.Description, in.Status, in.Label, in.IsPublished, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskResponse{Nanoid: nanoid}, nil
}
