CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    email text NOT NULL UNIQUE,
    hashed_password text NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp
)