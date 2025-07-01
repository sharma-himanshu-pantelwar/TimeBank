-- Users
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
username VARCHAR(30) NOT NULL,
email TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
location VARCHAR(20) NOT NULL,
availability BOOLEAN DEFAULT false,
available_credits INT NOT NULL DEFAULT 0,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Auth Sessions
CREATE TABLE IF NOT EXISTS auth_sessions (
session_id UUID PRIMARY KEY,
user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash TEXT NOT NULL,
expires_at TIMESTAMPTZ NOT NULL,
issued_at TIMESTAMPTZ NOT NULL,
UNIQUE(user_id) --having this wouldn't allow multiple sessions for a user
);

-- Skills 
-- CREATE TYPE IF NOT EXISTS skill_status_types AS ENUM('inactive','active');
-- CREATE TYPE skill_service_types AS ENUM('needed','offered');
CREATE TABLE IF NOT EXISTS skills(
    skill_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    description TEXT,
    skill_status skill_status_types DEFAULT 'inactive',
    min_time_required INT NOT NULL,
    UNIQUE(user_id,name)
);


-- Time Credits
CREATE TABLE IF NOT EXISTS time_credits (
id UUID PRIMARY KEY,
earned_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
spent_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
value INT NOT NULL,
transaction_at TIMESTAMPTZ NOT NULL
);

-- Requests
CREATE TABLE IF NOT EXISTS requestS(
    req_id SERIAL PRIMARY KEY,
    from_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id INT NOT NULL REFERENCES skills(skill_id) ON DELETE CASCADE,
    accepted BOOLEAN DEFAULT false,
    requested_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- TimeBank Service Sessions
-- CREATE TYPE session_status_types AS ENUM('created','started','completed','cancelled');
CREATE TABLE IF NOT EXISTS helping_sessions (
id SERIAL PRIMARY KEY,
sender_id INT REFERENCES users(id) ON DELETE SET NULL,
receiver_id INT REFERENCES users(id) ON DELETE SET NULL,
skill_shared_id INT REFERENCES skills(skill_id),
time_taken NUMERIC(4,2) CHECK(time_taken >= 0),
-- session_status session_status_types DEFAULT 'created',
-- scheduled_at TIMESTAMPTZ,
started_at TIMESTAMPTZ,
completed_at TIMESTAMPTZ,
-- created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Feedback
CREATE TABLE IF NOT EXISTS feedback (
id SERIAL PRIMARY KEY,
session_id INT UNIQUE REFERENCES helping_sessions(id) ON DELETE CASCADE,
rater_id INT REFERENCES users(id) ON DELETE CASCADE,
ratee_id INT REFERENCES users(id) ON DELETE CASCADE,
rating INT CHECK(rating BETWEEN 1 AND 5),
comments TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
UNIQUE(session_id, rater_id)
);


-- -- Function: Prevent Negative Credits
-- CREATE OR REPLACE FUNCTION prevent_negative_credits() 
-- RETURNS TRIGGER AS $body$
-- BEGIN
--     IF NEW.available_credits < 0 THEN
--         RAISE EXCEPTION 'Insufficient credits';
--     END IF;
--     RETURN NEW;
-- END;
-- $body$ 
-- LANGUAGE plpgsql;

-- -- Trigger to enforce the check
-- CREATE TRIGGER check_balance
-- BEFORE UPDATE ON users
-- FOR EACH ROW
-- EXECUTE FUNCTION prevent_negative_credits();
