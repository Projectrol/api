CREATE TABLE IF NOT EXISTS projects_members (
    id serial primary key,
    member_id serial NOT NULL,
    project_id serial NOT NULL,
    created_at timestamp DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    UNIQUE(project_id, member_id),
    FOREIGN KEY (member_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
)