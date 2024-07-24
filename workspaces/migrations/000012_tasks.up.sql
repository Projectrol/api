CREATE TABLE IF NOT EXISTS tasks (
    id serial primary key,
    nanoid text NOT NULL UNIQUE,
    project_id serial NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    status text NOT NULL,
    label text NOT NULL,
    is_published boolean,
    created_at timestamp DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    created_by serial NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_nanoid_tasks ON tasks(nanoid);