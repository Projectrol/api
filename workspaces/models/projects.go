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

type SearchProjectsInput struct {
	Q string `json:"q"`
}

type CreateProjectDocumentInput struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type UpdateProjectDocumentInput struct {
	Name    string `json:"name"`
	Content string `json:"content"`
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
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
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
	searchValue := ""
	if in.Q != "*" && in.Q != "*," {
		searchValue = in.Q
	}
	log.Print(searchValue)
	queryStr := fmt.Sprintf(`SELECT P.id, P.workspace_id, P.created_by, P.name, P.slug, P.description, P.dtstart, P.dtend, P.created_at 
							FROM projects_members PM
							LEFT JOIN projects P ON PM.project_id = P.id
							WHERE (workspace_id = %d AND PM.member_id = %d AND P.is_private = true AND name ILIKE %s) 
							OR (workspace_id = %d AND P.is_private = false AND name ILIKE %s)`,
		in.WorkspaceId, in.UserId, "'%"+searchValue+"%'", in.WorkspaceId, "'%"+searchValue+"%'")
	rows, err := m.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := &pb.Project{}
		err = rows.Scan(&project.Id, &project.WorkspaceId, &project.CreatedBy,
			&project.Name, &project.Slug, &project.Description, &project.Dtstart, &project.Dtend, &project.CreatedAt)
		if err == nil {
			projects = append(projects, project)
		} else {
			log.Print(err.Error())
		}
	}
	log.Print(len(projects))
	return &pb.GetProjectsByWorkspaceIdResponse{Projects: projects}, nil
}

func (m *ProjectsModel) GetProjectDetails(ctx context.Context, in *pb.GetProjectDetailsRequest) (*pb.GetProjectDetailsResponse, error) {
	details := &pb.ProjectDetails{}
	project := &pb.Project{}
	err := m.DB.QueryRow(`SELECT id, workspace_id, created_by, name, slug, summary, description, dtstart, dtend, created_at 
					FROM projects WHERE slug = $1`, in.ProjectSlug).
		Scan(&project.Id, &project.WorkspaceId, &project.CreatedBy,
			&project.Name, &project.Slug, &project.Summary,
			&project.Description, &project.Dtstart, &project.Dtend, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	details.Project = project
	rows, err := m.DB.Query("SELECT member_id from projects_members WHERE project_id = $1", project.Id)
	var memberIds []int32
	if err == nil {
		for rows.Next() {
			var memberId int32
			err = rows.Scan(&memberId)
			if err == nil {
				memberIds = append(memberIds, memberId)
			}
		}
		details.MemberIds = memberIds
	}
	return &pb.GetProjectDetailsResponse{Details: details}, nil
}

func (m *ProjectsModel) CreateProjectDocument(ctx context.Context, in *pb.CreateProjectDocumentRequest) (*pb.CreateProjectDocumentResponse, error) {
	var projectId int
	err := m.DB.QueryRow("SELECT id FROM projects WHERE slug = $1", in.ProjectSlug).Scan(&projectId)
	if err != nil {
		return nil, err
	}

	nanoid := common.GenerateNanoid(10)
	_, err = m.DB.Exec(`INSERT INTO project_documents (created_by, updated_by, project_id, nanoid, name, content) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		in.UserId, in.UserId, projectId, nanoid, in.Name, in.Content)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProjectDocumentResponse{Nanoid: nanoid}, nil
}

func (m *ProjectsModel) GetProjectDocuments(ctx context.Context, in *pb.GetProjectDocumentsRequest) (*pb.GetProjectDocumentsResponse, error) {
	var projectId int
	err := m.DB.QueryRow("SELECT id FROM projects WHERE slug = $1", in.ProjectSlug).Scan(&projectId)
	if err != nil {
		return nil, err
	}
	rows, err := m.DB.Query("SELECT created_by, updated_by, nanoid, name, created_at, updated_at FROM project_documents WHERE project_id = $1 ORDER BY created_at DESC", projectId)
	if err != nil {
		return nil, err
	}
	var documents []*pb.ProjectDocument
	for rows.Next() {
		document := &pb.ProjectDocument{}
		err := rows.Scan(&document.CreatedBy, &document.UpdatedBy, &document.Nanoid, &document.Name, &document.CreatedAt, &document.UpdatedAt)
		if err == nil {
			documents = append(documents, document)
		}
	}
	return &pb.GetProjectDocumentsResponse{Documents: documents}, nil
}

func (m *ProjectsModel) GetProjectDocumentDetails(ctx context.Context, in *pb.GetProjectDocumentDetailsRequest) (*pb.GetProjectDocumentDetailsResponse, error) {
	details := &pb.ProjectDocumentDetails{}
	err := m.DB.QueryRow(`SELECT created_by, updated_by, nanoid, name, content, created_at, updated_at 
						FROM project_documents WHERE nanoid = $1`, in.Nanoid).
		Scan(&details.CreatedBy, &details.UpdatedBy, &details.Nanoid, &details.Name, &details.Content, &details.CreatedAt, &details.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &pb.GetProjectDocumentDetailsResponse{
		Details: details,
	}, nil
}

func (m *ProjectsModel) UpdateProjectDocumentDetails(ctx context.Context, in *pb.UpdateProjectDocumentDetailsRequest) (*pb.UpdateProjectDocumentDetailsResponse, error) {
	_, err := m.DB.Exec(`UPDATE project_documents SET name = $1, content = $2, updated_by = $3, updated_at = NOW() at time zone 'utc' WHERE nanoid = $4`,
		in.Name, in.Content, in.UpdatedBy, in.Nanoid)
	if err != nil {
		return nil, err
	}
	response, err := m.GetProjectDocumentDetails(ctx, &pb.GetProjectDocumentDetailsRequest{Nanoid: in.Nanoid})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProjectDocumentDetailsResponse{
		Details: response.Details,
	}, nil
}
