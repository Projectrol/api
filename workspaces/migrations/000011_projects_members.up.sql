CREATE TABLE IF NOT EXISTS projects_members (
    id serial primary key,
    member_id serial NOT NULL,
    created_at timestamp DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (member_id) REFERENCES users(id)
)