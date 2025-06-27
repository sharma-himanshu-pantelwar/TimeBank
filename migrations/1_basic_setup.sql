-- Users
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
username VARCHAR(30) NOT NULL,
email TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
balance_hours NUMERIC(5,2) DEFAULT 0.0 CHECK(balance_hours >= 0),
reputation NUMERIC(3,2) DEFAULT 5.0 CHECK(reputation BETWEEN 0 and 5),
skills_offered TEXT[] DEFAULT '{}',
skills_needed TEXT[] DEFAULT '{}',
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Auth Sessions
CREATE TABLE IF NOT EXISTS auth_sessions (
id UUID PRIMARY KEY,
user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash TEXT NOT NULL,
expires_at TIMESTAMPTZ NOT NULL,
issued_at TIMESTAMPTZ NOT NULL
-- UNIQUE(user_id) --having this wouldn't allow multiple sessions for a user
);

-- TimeBank Service Sessions
CREATE TYPE session_status AS ENUM('created','started','completed','cancelled');
CREATE TABLE IF NOT EXISTS time_sessions (
id SERIAL PRIMARY KEY,
helper_id INT REFERENCES users(id) ON DELETE SET NULL,
recipient_id INT REFERENCES users(id) ON DELETE SET NULL,
skill VARCHAR(50) NOT NULL,
hours_estimate NUMERIC(4,2) NOT NULL CHECK(hours_estimate > 0),
actual_hours NUMERIC(4,2) CHECK(actual_hours >= 0),
status session_status DEFAULT 'created',
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
session_id INT UNIQUE REFERENCES time_sessions(id) ON DELETE CASCADE,
rater_id INT REFERENCES users(id) ON DELETE CASCADE,
ratee_id INT REFERENCES users(id) ON DELETE CASCADE,
rating INT CHECK(rating BETWEEN 1 AND 5),
comments TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
UNIQUE(session_id, rater_id)
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
RETURNS TRIGGER AS $func$
BEGIN
    IF NEW.balance_hours < 0 THEN
        RAISE EXCEPTION 'Insufficient balance';
    END IF;
    RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

-- Trigger to use the function

CREATE TRIGGER check_balance
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_negative_balance();