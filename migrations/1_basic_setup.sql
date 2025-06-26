-- Users
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
name VARCHAR(30) NOT NULL,
email TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
balance_hours NUMERIC(5,2) DEFAULT 0.0 CHECK(balance_hours >= 0),
reputation NUMERIC(3,2) DEFAULT 5.0,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Auth Sessions
CREATE TABLE IF NOT EXISTS sessions (
id UUID PRIMARY KEY,
user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash TEXT NOT NULL,
expires_at TIMESTAMPTZ NOT NULL,
issued_at TIMESTAMPTZ NOT NULL,
UNIQUE(user_id)
);

-- User Skills
CREATE TABLE IF NOT EXISTS user_skills (
id SERIAL PRIMARY KEY,
user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
skill VARCHAR(50) NOT NULL,
type VARCHAR(10) CHECK(type IN ('offered', 'needed')),
UNIQUE(user_id, skill, type)
);

-- TimeBank Service Sessions
CREATE TABLE IF NOT EXISTS time_sessions (
id SERIAL PRIMARY KEY,
helper_id INT REFERENCES users(id) ON DELETE SET NULL,
recipient_id INT REFERENCES users(id) ON DELETE SET NULL,
skill VARCHAR(50) NOT NULL,
hours_estimate NUMERIC(4,2) NOT NULL CHECK(hours_estimate > 0),
actual_hours NUMERIC(4,2) CHECK(actual_hours > 0),
status VARCHAR(20) DEFAULT 'created' CHECK (status IN ('created', 'started', 'completed', 'cancelled')),
notes TEXT,
completed_notes TEXT,
scheduled_at TIMESTAMPTZ,
started_at TIMESTAMPTZ,
completed_at TIMESTAMPTZ,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Feedback
CREATE TABLE IF NOT EXISTS feedback (
id SERIAL PRIMARY KEY,
session_id INT REFERENCES time_sessions(id) ON DELETE CASCADE,
from_user_id INT REFERENCES users(id),
to_user_id INT REFERENCES users(id),
rating INT CHECK(rating BETWEEN 1 AND 5),
comments TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
UNIQUE(session_id, from_user_id)
);

-- Time Credit Ledger
CREATE TABLE IF NOT EXISTS time_credits (
id SERIAL PRIMARY KEY,
session_id INT REFERENCES time_sessions(id),
from_user_id INT REFERENCES users(id),
to_user_id INT REFERENCES users(id),
hours NUMERIC(4,2) NOT NULL CHECK (hours > 0),
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Function: Prevent Negative Balance
CREATE OR REPLACE FUNCTION prevent_negative_balance()
RETURNS TRIGGER AS $$
BEGIN
IF NEW.balance_hours < 0 THEN
RAISE EXCEPTION 'Insufficient balance';
END IF;
RETURN NEW;
END;