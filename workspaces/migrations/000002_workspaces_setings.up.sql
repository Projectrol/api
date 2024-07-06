CREATE TABLE IF NOT EXISTS workspaces_settings (
    id serial primary key,
    workspace_id serial UNIQUE NOT NULL,
    logo text NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
)