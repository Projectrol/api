CREATE TABLE IF NOT EXISTS users_settings (
    id serial primary key,
    user_id serial NOT NULL UNIQUE,
    name text,
    avatar text,
    theme text,
    phone_no text,
    created_at timestamp DEFAULT (NOW() at time zone 'utc'),
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id)
)