CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    email text UNIQUE NOT NULL,
    hashed_password text NOT NULL,
    created_at timestamp default (now() at time zone 'utc'),
    updated_at timestamp
)