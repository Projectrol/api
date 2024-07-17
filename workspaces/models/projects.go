package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	IsPrivate   bool   `json:"is_private"`
}

func NewProjectsModel(DB *sql.DB) *ProjectsModel {
	return &ProjectsModel{
		DB: DB,
	}
}

func (m *ProjectsModel) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	slug := common.GenerateSlugName(in.Name) + "-" + common.GenerateNanoid(8)
	dtstart, _ := strconv.ParseInt(fmt.Sprintf("%d", in.Dtstart), 10, 64)
	dtstartTime := time.Unix(dtstart, 0).UTC()
	dtEnd, _ := strconv.ParseInt(fmt.Sprintf("%d", in.Dtend), 10, 64)
	dtEndTime := time.Unix(dtEnd, 0).UTC()
	var id int32
	err := m.DB.QueryRow(`INSERT INTO projects(workspace_id, created_by, slug, name, summary, description, dtstart, dtend, is_private) 
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ID`,
		in.WorkspaceId,
		in.CreatedBy,
		slug,
		in.Name,
		in.Summary,
		in.Description,
		dtstartTime,
		dtEndTime,
		in.IsPrivate,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	m.DB.Exec(`INSERT INTO projects_members(member_id, project_id) VALUES($1, $2)`,
		in.CreatedBy,
		id,
	)
	return &pb.CreateProjectResponse{Id: id}, nil
}

func (m *ProjectsModel) GetProjectsByWorkspaceId(ctx context.Context, in *pb.GetProjectsByWorkspaceIdRequest) (*pb.GetProjectsByWorkspaceIdResponse, error) {
	var projects []*pb.Project

	rows, err := m.DB.Query(`SELECT P.id, P.workspace_id, P.created_by, P.name, P.slug, P.description, P.dtstart, P.dtend, P.created_at 
							FROM projects_members PM
							LEFT JOIN projects P ON PM.project_id = P.id
							WHERE (workspace_id = $1 AND PM.member_id = $2 AND P.is_private = true) 
							OR (workspace_id = $1 AND P.is_private = false)`, in.WorkspaceId, in.UserId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := &pb.Project{}
		err = rows.Scan(&project.Id, &project.WorkspaceId, &project.CreatedBy,
			&project.Name, &project.Slug, &project.Description, &project.Dtstart, &project.Dtend, &project.CreatedAt)
		if err == nil {
			log.Print(2)
			projects = append(projects, project)
		} else {
			log.Print(err.Error())
		}
	}
	return &pb.GetProjectsByWorkspaceIdResponse{Projects: projects}, nil
}

func (m *ProjectsModel) GetProjectDetails(ctx context.Context, in *pb.GetProjectDetailsRequest) (*pb.GetProjectDetailsResponse, error) {
	project := &pb.Project{}
	err := m.DB.QueryRow(`SELECT id, workspace_id, created_by, name, slug, summary, description, dtstart, dtend, created_at 
					FROM projects WHERE slug = $1`, in.ProjectSlug).
		Scan(&project.Id, &project.WorkspaceId, &project.CreatedBy,
			&project.Name, &project.Slug, &project.Summary,
			&project.Description, &project.Dtstart, &project.Dtend, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pb.GetProjectDetailsResponse{Project: project}, nil
}
