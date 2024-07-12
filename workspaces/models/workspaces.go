package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	nanoid := common.GenerateNanoid(10)
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

	for i := 0; i < 2; i++ {
		var roleId int
		var permissionIds []int
		roleName := ""
		switch i {
		case 0:
			roleName = "Owner"
			rows, err := m.DB.Query("SELECT id FROM permissions ORDER BY created_at ASC")
			if err == nil {
				for rows.Next() {
					var permissionId int
					err = rows.Scan(&permissionId)
					if err == nil {
						permissionIds = append(permissionIds, permissionId)
					}
				}
			}
		case 1:
			roleName = "Admin"
			rows, err := m.DB.Query("SELECT id FROM permissions WHERE resource_tag != 'workspaces' ORDER BY created_at ASC")
			if err == nil {
				for rows.Next() {
					var permissionId int
					err = rows.Scan(&permissionId)
					if err == nil {
						permissionIds = append(permissionIds, permissionId)
					}
				}
			}
		}
		err = m.DB.QueryRow("INSERT INTO workspace_roles(workspace_id, role_name) VALUES($1, $2) RETURNING id", id, roleName).
			Scan(&roleId)
		if err != nil {
			log.Print("error at workspace_roles. Error:" + err.Error())
		}
		if err == nil {
			if roleName == "Owner" {
				_, err = m.DB.Exec("INSERT INTO workspace_members(workspace_id, user_id, role_id) VALUES($1, $2, $3)", id, ownerId, roleId)
				if err != nil {
					log.Print("error at role_permissions. Error:" + err.Error())
				}
			}

			for _, permissionId := range permissionIds {
				_, err = m.DB.Exec("INSERT INTO role_permissions(role_id, permission_id) VALUES($1, $2)", roleId, permissionId)
				if err != nil {
					log.Print("error at role_permissions. Error:" + err.Error())
				}
			}
		}
	}

	if err != nil {
		return "", err
	}

	return nanoid, nil
}

func (m *WorkspaceModel) GetWorkspacesByUserId(ctx context.Context, userId int32) (*pb.GetWorkspacesByUserIdResponse, error) {
	rows, err := m.DB.Query(`SELECT WS.id, nanoid, name, slug FROM workspaces WS 
							LEFT JOIN workspace_members WM ON WS.id = WM.workspace_id
							WHERE WM.user_id = $1`, userId)
	if err != nil {
		return &pb.GetWorkspacesByUserIdResponse{Workspaces: make([]*pb.Workspace, 0)}, err
	}
	var workspaces []*pb.Workspace
	for rows.Next() {
		w := &pb.Workspace{}
		err = rows.Scan(&w.Id, &w.Nanoid, &w.Name, &w.Slug)
		if err == nil {
			workspaces = append(workspaces, w)
		}
	}
	return &pb.GetWorkspacesByUserIdResponse{Workspaces: workspaces}, nil
}

func (m *WorkspaceModel) GetWorkspaceDetails(ctx context.Context, in *pb.GetWorkspaceDetailsRequest) (*pb.GetWorkspaceDetailsResponse, error) {
	details := &pb.GetWorkspaceDetailsResponse{}
	var id int32
	var nanoid string
	var name string
	var slug string
	err := m.DB.QueryRow("SELECT id, nanoid, name, slug FROM workspaces WHERE id = $1", in.Id).
		Scan(&id,
			&nanoid,
			&name,
			&slug)
	if err != nil {
		log.Print("148")
		return nil, err
	}

	details.GeneralInformation = &pb.Workspace{
		Id:     id,
		Nanoid: nanoid,
		Name:   name,
		Slug:   slug,
	}

	var logo string
	err = m.DB.QueryRow("SELECT logo FROM workspaces_settings WHERE workspace_id=$1", in.Id).Scan(&logo)
	if err == nil {
		details.Settings = &pb.WorkspaceSettings{Logo: logo}
	}

	return details, nil
}

func (m *WorkspaceModel) GetPermissions(ctx context.Context, in *pb.EmptyRequest) (*pb.GetPermissionsResponse, error) {
	rows, err := m.DB.Query(`SELECT id, resource_tag, title, description, 
							can_read, can_create, can_update, can_delete FROM permissions 
							ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	var permissions []*pb.Permission
	for rows.Next() {
		permission := &pb.Permission{}
		err = rows.Scan(&permission.Id,
			&permission.ResourceTag,
			&permission.Title,
			&permission.Description,
			&permission.CanRead, &permission.CanCreate, &permission.CanUpdate, &permission.CanDelete,
		)
		if err == nil {
			permissions = append(permissions, permission)
		}
	}
	return &pb.GetPermissionsResponse{Permissions: permissions}, nil
}

func (m *WorkspaceModel) UpdateRolePermission(ctx context.Context, in *pb.UpdateRolePermissionRequest) (*pb.UpdateRolePermissionResponse, error) {
	roleId := in.RoleId
	rTag := in.ResourceTag
	action := in.Action
	updateType := in.UpdateType

	var permissionId int
	queryStr := fmt.Sprintf(`SELECT id FROM permissions WHERE resource_tag = '%s' AND %s = 'true'`, rTag, action)
	log.Print(queryStr)
	err := m.DB.QueryRow(queryStr).
		Scan(&permissionId)
	if err != nil {
		return nil, err
	}

	switch updateType {
	case "add":
		_, err = m.DB.Exec("INSERT INTO role_permissions(role_id, permission_id) VALUES($1, $2)", roleId, permissionId)
	case "remove":
		_, err = m.DB.Exec("DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2", roleId, permissionId)
	}

	if err != nil {
		return nil, err
	}

	return &pb.UpdateRolePermissionResponse{}, nil
}

func (m *WorkspaceModel) GetWorkspaceRoles(ctx context.Context, in *pb.GetWorkspaceRolesRequest) (*pb.GetWorkspaceRolesResponse, error) {
	rows, err := m.DB.Query("SELECT id, workspace_id, role_name FROM workspace_roles WHERE workspace_id = $1", in.Id)
	if err != nil {
		return nil, err
	}
	var roles []*pb.WorkspaceRole
	for rows.Next() {
		role := &pb.WorkspaceRole{}
		var roleId int32
		var workspaceId int32
		var roleName string
		err = rows.Scan(&roleId, &workspaceId, &roleName)
		if err == nil {
			role.Id = roleId
			role.WorkspaceId = workspaceId
			role.RoleName = roleName
			pRows, err := m.DB.Query(`SELECT permissions.id, resource_tag, title, description, can_read, can_create, can_update, can_delete 
											FROM permissions 
											LEFT JOIN role_permissions
											ON permissions.id = role_permissions.permission_id
											WHERE role_permissions.role_id = $1`, role.Id)
			if err != nil {
				role.Permissions = make([]*pb.Permission, 0)
				continue
			} else {
				var permissions []*pb.Permission
				for pRows.Next() {
					permission := &pb.Permission{}
					err = pRows.Scan(&permission.Id,
						&permission.ResourceTag,
						&permission.Title,
						&permission.Description,
						&permission.CanRead, &permission.CanCreate, &permission.CanUpdate, &permission.CanDelete,
					)
					if err == nil {
						permissions = append(permissions, permission)
					}
				}

				role.Permissions = permissions
			}
		}
		roles = append(roles, role)
	}

	return &pb.GetWorkspaceRolesResponse{Roles: roles}, nil
}

func (m *WorkspaceModel) GetUserRoleInWorkspace(ctx context.Context, in *pb.GetUserRoleInWorkspaceRequest) (*pb.GetUserRoleInWorkspaceResponse, error) {
	role := &pb.WorkspaceRole{}
	err := m.DB.QueryRow(`SELECT WR.id, WR.workspace_id, WR.role_name FROM workspace_roles WR
							LEFT JOIN workspace_members WM
							ON WR.id = WM.role_id
							WHERE WM.user_id = $1 AND WM.workspace_id = $2`, in.UserId, in.WorkspaceId).
		Scan(&role.Id, &role.WorkspaceId, &role.RoleName)
	if err != nil {
		return nil, err
	}
	pRows, err := m.DB.Query(`SELECT permissions.id, resource_tag, title, description, can_read, can_create, can_update, can_delete 
											FROM permissions 
											LEFT JOIN role_permissions
											ON permissions.id = role_permissions.permission_id
											WHERE role_permissions.role_id = $1`, role.Id)
	if err != nil {
		return nil, err
	}

	var permissions []*pb.Permission
	for pRows.Next() {
		permission := &pb.Permission{}
		err = pRows.Scan(&permission.Id,
			&permission.ResourceTag,
			&permission.Title,
			&permission.Description,
			&permission.CanRead, &permission.CanCreate, &permission.CanUpdate, &permission.CanDelete,
		)
		if err == nil {
			permissions = append(permissions, permission)
		}
	}
	role.Permissions = permissions

	return &pb.GetUserRoleInWorkspaceResponse{Role: role}, nil
}

func (m *WorkspaceModel) GetRoleIdOfUserWorkspaces(ctx context.Context, in *pb.GetRoleIdOfUserWorkspacesRequest) (*pb.GetRoleIdOfUserWorkspacesResponse, error) {
	rows, err := m.DB.Query("SELECT role_id, workspace_id FROM workspace_members WHERE user_id = $1", in.UserId)
	if err != nil {
		return nil, err
	}
	var workspaceRoleIdList []*pb.WorkspaceRoleId
	for rows.Next() {
		workspaceRoleId := &pb.WorkspaceRoleId{}
		err = rows.Scan(&workspaceRoleId.RoleId, &workspaceRoleId.WorkspaceId)
		if err == nil {
			workspaceRoleIdList = append(workspaceRoleIdList, workspaceRoleId)
		}
	}
	return &pb.GetRoleIdOfUserWorkspacesResponse{WorkpaceRoleIdList: workspaceRoleIdList}, nil
}

func (m *WorkspaceModel) CheckRoleValidForResource(ctx context.Context, in *pb.CheckRoleValidForResourceRequest) (*pb.CheckRoleValidForResourceResponse, error) {
	permissionType := ""
	roleId := in.RoleId
	resourceTag := in.ResourceTag

	switch in.Method {
	case "POST":
		permissionType = "can_create"
	case "PATCH":
		permissionType = "can_update"
	case "DELETE":
		permissionType = "can_delete"
	case "GET":
		permissionType = "can_read"
	}

	queryStr := fmt.Sprintf(`SELECT COUNT(*) FROM permissions 
							LEFT JOIN role_permissions 
							ON permissions.id = role_permissions.permission_id 
							WHERE role_permissions.role_id = %d AND %s = 'true' AND resource_tag = '%s'`, roleId, permissionType, resourceTag)
	var count int
	log.Print(queryStr)
	err := m.DB.QueryRow(queryStr).Scan(&count)
	if err != nil || count == 0 {
		return &pb.CheckRoleValidForResourceResponse{IsValid: false}, nil
	}
	return &pb.CheckRoleValidForResourceResponse{IsValid: true}, nil
}

func (m *WorkspaceModel) CreateNewRole(ctx context.Context, in *pb.CreateNewRoleRequest) (*pb.CreateNewRoleResponse, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM workspace_roles WHERE role_name ILIKE $1 AND workspace_id = $2", in.RoleName, in.WorkspaceId).
		Scan(&count)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("duplicate role_name")
	}

	var roleId int
	err = m.DB.QueryRow("INSERT INTO workspace_roles(workspace_id, role_name) VALUES ($1, $2) RETURNING id", in.WorkspaceId, in.RoleName).
		Scan(&roleId)

	if err != nil {
		return nil, err
	}

	for _, permissionId := range in.PermissionIds {
		m.DB.Exec("INSERT INTO role_permissions(role_id, permission_id) VALUES($1, $2)", roleId, permissionId)
	}

	return &pb.CreateNewRoleResponse{Id: int32(roleId)}, nil
}
