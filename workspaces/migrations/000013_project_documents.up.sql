CREATE TABLE IF NOT EXISTS project_documents (
    id serial primary key,
    created_by serial NOT NULL,
    updated_by serial NOT NULL,
    project_id serial NOT NULL,
    nanoid text UNIQUE NOT NULL,
    name text NOT NULL,
    content text NOT NULL,
    created_at timestamp DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
)