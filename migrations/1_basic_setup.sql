-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    location VARCHAR(20) NOT NULL,
    availability BOOLEAN DEFAULT true,
    available_credits NUMERIC(5,2) NOT NULL DEFAULT 0.0,
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
-- CREATE TYPE skill_status_types AS ENUM('inactive','active');
-- CREATE TYPE skill_service_types AS ENUM('needed','offered');
CREATE TABLE IF NOT EXISTS skills(
    skill_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    description TEXT,
    skill_status skill_status_types DEFAULT 'inactive',
    min_time_required NUMERIC(4,2) NOT NULL DEFAULT 0.0,
    UNIQUE(user_id,name)
);


-- Time Credits
CREATE TABLE IF NOT EXISTS time_credits (
id SERIAL PRIMARY KEY,
earned_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
spent_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
value NUMERIC(5,2) NOT NULL,
transaction_at TIMESTAMPTZ NOT NULL
);

-- Requests
CREATE TABLE IF NOT EXISTS requests(
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
sender_id INT NOT NULL REFERENCES users(id) ,
receiver_id INT NOT NULL REFERENCES users(id),
skill_shared_id INT NOT NULL REFERENCES skills(skill_id),
time_taken NUMERIC(5,2) CHECK(time_taken >= 0),
-- session_status session_status_types DEFAULT 'created',
-- scheduled_at TIMESTAMPTZ,
started_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
completed_at TIMESTAMPTZ
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


-- Function: Prevent Negative Credits
CREATE OR REPLACE FUNCTION prevent_negative_credits() 
RETURNS TRIGGER AS $body$
BEGIN
    IF NEW.available_credits < 0 THEN
        RAISE EXCEPTION 'Insufficient credits';
    END IF;
    RETURN NEW;
END;
$body$ 
LANGUAGE plpgsql;

-- Trigger to enforce the check
CREATE TRIGGER check_balance
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_negative_credits();


CREATE OR REPLACE FUNCTION calculate_time_taken()
RETURNS TRIGGER AS $$
BEGIN
  -- Only calculate if completed_at is being updated
  IF NEW.completed_at IS NOT NULL AND OLD.completed_at IS DISTINCT FROM NEW.completed_at THEN
    -- time_taken in hours
    NEW.time_taken := ROUND(EXTRACT(EPOCH FROM NEW.completed_at - NEW.started_at) / 60, 2);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_calculate_time_taken
BEFORE UPDATE ON helping_sessions
FOR EACH ROW
EXECUTE FUNCTION calculate_time_taken();
