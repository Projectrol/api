CREATE TABLE IF NOT EXISTS role_permissions (
    id serial primary key,
    role_id serial NOT NULL,
    permission_id serial NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    UNIQUE (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES workspace_roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
)