CREATE TABLE IF NOT EXISTS workspace_roles (
    id serial primary key,
    workspace_id serial NOT NULL,
    role_name text NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    UNIQUE (workspace_id, role_name),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
)