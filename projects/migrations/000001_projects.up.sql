CREATE TABLE IF NOT EXISTS projects (
    id serial primary key,
    workspace_id serial NOT NULL,
    slug text NOT NULL,
    name text NOT NULL,
    summary text NOT NULL,
    description text NOT NULL,
    dtstart timestamp,
    dtend timestamp,
    created_by serial NOT NULL,
    created_at timestamp NOT NULL DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    UNIQUE (workspace_id, slug),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
)