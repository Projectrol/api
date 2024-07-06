package models

import (
	"context"
	"database/sql"
	"errors"

	common "github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type Workspace struct {
	Id        int            `json:"id"`
	Nanoid    string         `json:"nanoid"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	OwnerId   int            `json:"owner_id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt sql.NullString `json:"updated_at"`
}

type WorkspaceSettings struct {
	Id          int            `json:"id"`
	WorkspaceId int            `json:"workspace_id"`
	Logo        string         `json:"logo"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   sql.NullString `json:"updated_at"`
}

type CreateWorkspaceInput struct {
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type WorkspaceModel struct {
	DB *sql.DB
}

func NewWorkspaceModel(DB *sql.DB) *WorkspaceModel {
	return &WorkspaceModel{DB}
}

func (m *WorkspaceModel) Insert(ctx context.Context, input *pb.CreateWorkspaceRequest) (string, error) {
	slug := common.GenerateSlugName(input.Name)
	nanoid := common.GenerateNanoid()
	name := input.Name
	ownerId := input.OwnerId
	logo := input.Logo
	var existedName string
	row := m.DB.QueryRow("SELECT name FROM workspaces WHERE name=$1", name)
	err := row.Scan(&existedName)
	if err == nil {
		return "", errors.New("duplicated name")
	}
	_, err = m.DB.Exec("INSERT INTO workspaces(nanoid, name, slug, owner_id) VALUES($1, $2, $3, $4)", nanoid, name, slug, ownerId)
	if err != nil {
		return "", err
	}
	var id int
	err = m.DB.QueryRow("SELECT id from workspaces WHERE nanoid=$1", nanoid).Scan(&id)
	if err != nil {
		return "", err
	}

	_, err = m.DB.Exec("INSERT INTO workspaces_settings(workspace_id, logo) VALUES($1, $2)", id, logo)

	if err != nil {
		return "", err
	}

	return nanoid, nil
}

func (m *WorkspaceModel) GetWorkspacesByUserId(ctx context.Context, userId int32) (*pb.GetWorkspacesByUserIdResponse, error) {
	rows, err := m.DB.Query("SELECT id, nanoid, name, slug FROM workspaces WHERE owner_id=$1", userId)
	if err != nil {
		return &pb.GetWorkspacesByUserIdResponse{Workspaces: make([]*pb.Workspace, 0)}, err
	}
	var workspaces []*pb.Workspace
	for rows.Next() {
		w := &pb.Workspace{}
		err = rows.Scan(&w.Id, &w.Nanoid, &w.Name, &w.Slug)
		if err == nil {
			ws := &pb.WorkspaceSettings{}
			err = m.DB.QueryRow("SELECT logo FROM workspaces_settings WHERE workspace_id=$1", w.Id).Scan(&ws.Logo)
			if err == nil {
				w.Settings = ws
			}
			workspaces = append(workspaces, w)
		}
	}
	return &pb.GetWorkspacesByUserIdResponse{Workspaces: workspaces}, nil
}
