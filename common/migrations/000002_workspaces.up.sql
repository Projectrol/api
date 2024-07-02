CREATE TABLE IF NOT EXISTS workspaces (
    id serial primary key,
    nanoid text UNIQUE NOT NULL,
    slug text UNIQUE NOT NULL,
    name text UNIQUE NOT NULL,
    owner_id serial NOT NULL,
    created_at timestamp default (now() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (owner_id) REFERENCES users(id)
)