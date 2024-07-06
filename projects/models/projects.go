package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type ProjectsModel struct {
	DB *sql.DB
}

type CreateProjectInput struct {
	WorkspaceId int    `json:"workspace_id"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Dtstart     int64  `json:"dtstart"`
	Dtend       int64  `json:"dtend"`
}

func NewProjectsModel(DB *sql.DB) *ProjectsModel {
	return &ProjectsModel{
		DB: DB,
	}
}

func (m *ProjectsModel) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	slug := common.GenerateSlugName(in.Name)
	dtstart, _ := strconv.ParseInt(fmt.Sprintf("%d", in.Dtstart), 10, 64)
	dtstartTime := time.Unix(dtstart, 0).UTC()
	dtEnd, _ := strconv.ParseInt(fmt.Sprintf("%d", in.Dtend), 10, 64)
	dtEndTime := time.Unix(dtEnd, 0).UTC()
	var id int32
	err := m.DB.QueryRow(`INSERT INTO projects(workspace_id, created_by, slug, name, summary, description, dtstart, dtend) 
				VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID`,
		in.WorkspaceId,
		in.CreatedBy,
		slug,
		in.Name,
		in.Summary,
		in.Description,
		dtstartTime,
		dtEndTime,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProjectResponse{Id: id}, nil
}

func (m *ProjectsModel) GetProjectsByWorkspaceId(ctx context.Context, in *pb.GetProjectsByWorkspaceIdRequest) (*pb.GetProjectsByWorkspaceIdResponse, error) {
	var projects []*pb.Project

	rows, err := m.DB.Query(`SELECT id, workspace_id, created_by, name, slug, description, dtstart, dtend, created_at 
							FROM projects WHERE workspace_id = $1 `, in.WorkspaceId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := &pb.Project{}
		err = rows.Scan(&project.Id, &project.WorkspaceId, &project.CreatedBy,
			&project.Name, &project.Slug, &project.Description, &project.Dtstart, &project.Dtend, &project.CreatedAt)
		if err == nil {
			projects = append(projects, project)
		}
	}

	return &pb.GetProjectsByWorkspaceIdResponse{Projects: projects}, nil
}
