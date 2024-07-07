CREATE TABLE IF NOT EXISTS task_instances (
    id serial primary key,
    task_entity_id serial NOT NULL,
    dtstart text NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (task_entity_id) REFERENCES tasks(id)
)