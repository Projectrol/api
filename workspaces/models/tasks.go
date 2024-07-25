package models

import (
	"context"
	"database/sql"
	"log"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type TaskModel struct {
	DB *sql.DB
}

type CreateTaskInput struct {
	ProjectSlug string `json:"project_slug"`
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
	var projectId int
	err := m.DB.QueryRow("SELECT id FROM projects WHERE slug = $1", in.ProjectSlug).Scan(&projectId)
	if err != nil {
		return nil, err
	}
	_, err = m.DB.Exec(`INSERT INTO tasks (nanoid, project_id, title, description, status, label, is_published, created_by) 
					VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		nanoid, projectId, in.Title, in.Description, in.Status, in.Label, in.IsPublished, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskResponse{Nanoid: nanoid}, nil
}

func (m *TaskModel) GetProjectTasks(ctx context.Context, in *pb.GetProjectTasksRequest) (*pb.GetProjectTasksResponse, error) {
	var projectId int
	err := m.DB.QueryRow("SELECT id FROM projects WHERE slug = $1", in.ProjectSlug).Scan(&projectId)
	if err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(`SELECT nanoid, project_id, title, description, status, label, is_published, created_at 
							FROM tasks WHERE project_id = $1 AND is_published = true`, projectId)
	if err != nil {
		return nil, err
	}
	var tasks []*pb.Task
	for rows.Next() {
		task := &pb.Task{}
		err = rows.Scan(&task.Nanoid, &task.ProjectId, &task.Title, &task.Description,
			&task.Status, &task.Label, &task.IsPublished, &task.CreatedAt)
		if err == nil {
			tasks = append(tasks, task)
		} else {
			log.Print(err)
		}
	}
	return &pb.GetProjectTasksResponse{Tasks: tasks}, nil
}
