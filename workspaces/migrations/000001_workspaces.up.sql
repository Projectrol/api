CREATE TABLE IF NOT EXISTS workspaces (
    id serial primary key,
    nanoid text NOT NULL UNIQUE,
    name text NOT NULL UNIQUE,
    slug text NOT NULL UNIQUE,
    owner_id serial NOT NULL,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (owner_id) REFERENCES users(id)
)