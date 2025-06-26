CREATE TABLE IF NOT EXISTS users(
    uid serial PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions(
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(uid),
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    UNIQUE(user_id)
);