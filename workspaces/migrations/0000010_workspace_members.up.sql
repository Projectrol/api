CREATE TABLE IF NOT EXISTS workspace_members (
    id serial primary key,
    workspace_id serial NOT NULL,
    user_id serial NOT NULL,
    role_id serial NOT NULL,
    UNIQUE(workspace_id, user_id),
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES workspace_roles(id)
)