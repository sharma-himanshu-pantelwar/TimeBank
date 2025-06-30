package persistance

import (
	"fmt"
	"timebank/internal/core/session"
)

// SessionRepo is a struct which has db as pointer to Database
type SessionRepo struct {
	db *Database
}

// This function accepts d(database) as a parameter and assigns it to db in SessionRepo struct
func NewSessionRepo(d *Database) SessionRepo {
	return SessionRepo{db: d}
}

func (u *SessionRepo) CreateSession(session session.Session) error {
	_, err := u.db.db.Exec("INSERT INTO auth_sessions(session_id,user_id,token_hash,expires_at,issued_at) VALUES($1,$2,$3,$4,$5) ON CONFLICT (user_id) DO UPDATE SET session_id=EXCLUDED.session_id, token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at, issued_at = EXCLUDED.issued_at", session.Id, session.Uid, session.TokenHash, session.ExpiresAt, session.IssuedAt)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Session inserted into db : ", session.Id)
	return nil
}

func (u *SessionRepo) DeleteSession(uid int) error {
	query := "delete from sessions where user_id=$1"
	_, err := u.db.db.Query(query, uid)
	if err != nil {
		return err
	}
	return nil
}
func (u *SessionRepo) GetSession(id string) (session.Session, error) {
	var newSess session.Session
	query := "select session_id,user_id,token_hash,expires_at,issued_at from auth_sessions where session_id=$1"
	err := u.db.db.QueryRow(query, id).Scan(&newSess.Id, &newSess.Uid, &newSess.TokenHash, &newSess.ExpiresAt, &newSess.IssuedAt)
	if err != nil {
		return session.Session{}, err
	}
	return newSess, nil
}
