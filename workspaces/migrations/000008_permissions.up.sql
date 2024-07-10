CREATE TABLE IF NOT EXISTS permissions (
    id serial primary key,
    resource_tag text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    can_read boolean,
    can_create boolean,
    can_update boolean,
    can_delete boolean,
    UNIQUE(resource_tag, title),
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp
)