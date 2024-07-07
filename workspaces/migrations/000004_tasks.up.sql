CREATE TABLE IF NOT EXISTS tasks (
    id serial primary key,
    nanoid text NOT NULL UNIQUE,
    title text NOT NULL,
    description text NOT NULL,
    duration integer,
    type text NOT NULL,
    workspace_id serial NOT NULL,
    created_by serial NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
)