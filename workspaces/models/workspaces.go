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
	row = m.DB.QueryRow("SELECT id FROM workspaces WHERE nanoid=$1", nanoid)
	err = row.Scan(&id)
	if err != nil {
		return "", errors.New("something error")
	}
	return nanoid, nil
}
