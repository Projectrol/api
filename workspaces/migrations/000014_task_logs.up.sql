CREATE TABLE IF NOT EXISTS task_logs (
    id serial primary key,
    task_id serial NOT NULL,
    created_by serial NOT NULL,
    changed_field text NOT NULL,
    old_value text NOT NULL,
    new_value text NOT NULL,
    created_at timestamp DEFAULT(NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);