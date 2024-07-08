CREATE TABLE IF NOT EXISTS calendar_event_instances (
    id serial primary key,
    event_entity_id serial NOT NULL,
    dtstart text NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (event_entity_id) REFERENCES calendar_events(id)
)