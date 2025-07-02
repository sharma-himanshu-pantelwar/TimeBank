package persistance

import (
	"fmt"
	"time"
	helpsession "timebank/internal/core/help_session"
	"timebank/internal/core/skills"
	user "timebank/internal/core/user"
	hashpassword "timebank/pkg/hashPassword"
)

// * 1. Create User_Repo

// create a struct named UserRepo which stores a pointer to Database
type UserRepo struct {
	db *Database
}

// create a function NewUserRepo which accepts d (pointer to DB) as parameter and returns UserRepo
// which sets the accepted d as db for UserRepo
func NewUserRepo(d *Database) UserRepo {
	return UserRepo{db: d}
}

// Create user repo
func (u *UserRepo) CreateUser(newUser user.User) (user.User, error) {
	var id int
	query := "insert into users(username,email,password,location,availability,available_credits)values($1, $2, $3, $4, $5, $6) returning id"
	hashPass, err := hashpassword.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println("Error hashing password")
	}

	err = u.db.db.QueryRow(query, newUser.Username, newUser.Email, hashPass, newUser.Location, newUser.Availability, newUser.AvailableCredits).Scan(&id)
	if err != nil {
		return user.User{}, err
	}

	newUser.Id = id
	// send values to db
	return newUser, nil

}

func (u *UserRepo) GetUser(username string) (user.User, error) {
	var newUser user.User
	// fmt.Println(newUser)
	// u=>UserRepo      u.db=>Database type struct inside UserRepo   u.db.db=>Actual database inside the database struct(*sql.db)
	query := "select id,username,password,email,location,availability,available_credits from users where username=$1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Id, &newUser.Username, &newUser.Password, &newUser.Email, &newUser.Location, &newUser.Availability, &newUser.AvailableCredits)
	if err != nil {
		return user.User{}, err
	}
	// fmt.Println(newUser)
	return newUser, nil
}
func (u *UserRepo) GetUserByEmail(email string) (user.User, error) {
	var newUser user.User
	// fmt.Println(newUser)
	// u=>UserRepo      u.db=>Database type struct inside UserRepo   u.db.db=>Actual database inside the database struct(*sql.db)
	query := "select uid,username,password,email,location,availability,available_credits from users where email=$1"
	err := u.db.db.QueryRow(query, email).Scan(&newUser.Id, &newUser.Username, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	// fmt.Println(newUser)
	return newUser, nil
}
func (u *UserRepo) GetUserById(id int) (user.GetUserProfile, error) {
	var newUser user.GetUserProfile

	// fmt.Println(newUser)
	// u=>UserRepo      u.db=>Database type struct inside UserRepo   u.db.db=>Actual database inside the database struct(*sql.db)
	query := "select id,email,username,location,availability,available_credits from users where id=$1"
	err := u.db.db.QueryRow(query, id).Scan(&newUser.Uid, &newUser.Email, &newUser.Username, &newUser.Location, &newUser.Availability, &newUser.AvailableCredits)
	// resultUser.Email=newUser.Email
	// resultUser.Uid=newUser.Uid
	// resultUser.Username=newUser.Username
	if err != nil {
		return user.GetUserProfile{}, err
	}
	// fmt.Println(newUser)
	return newUser, nil
}

func (u *UserRepo) CreateSkill(userId int, newSkill skills.Skills) (skills.Skills, error) {
	var id int
	query := "insert into skills(user_id,name,description,skill_status,min_time_required)values($1, $2, $3, $4,$5) returning skill_id"

	// i need to get the user id from the authenticated user

	err := u.db.db.QueryRow(query, userId, newSkill.Name, newSkill.Description, newSkill.Status, newSkill.MinTimeRequired).Scan(&id)
	if err != nil {
		return skills.Skills{}, err
	}

	newSkill.Id = id
	newSkill.UserId = userId
	// send values to db
	return newSkill, nil

}

func (u *UserRepo) FindSkilledPerson(userId int, skillName string) ([]user.GetUsersWithSkills, error) {

	var people []user.GetUsersWithSkills
	query := "select users.id, users.username, users.email, skills.name,skills.skill_id, skills.description from users JOIN skills on users.id=skills.user_id where skills.name ILIKE $1 and users.id != $2;"

	rows, err := u.db.db.Query(query, "%"+skillName+"%", userId)

	if err != nil {
		fmt.Println("Error while running query        :             ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var person user.GetUsersWithSkills
		if err := rows.Scan(&person.Id, &person.Username, &person.Email, &person.SkillName, &person.SkillId, &person.SkillDescription); err != nil {
			fmt.Println("Error while scanning various rows")
			return nil, err
		}
		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error before returning people[]")
		return nil, err
	}
	// fmt.Println(people)  //empty array would go in case of no users found with that skill
	return people, nil

}

func (u *UserRepo) RenameSkill(userId int, newSkillName string, newSkillDescription string, skillId int) (skills.Skills, error) {

	var updatedSkill skills.Skills
	query := "update skills set name=$1, description=$2 where skills.skill_id= $3 and user_id=$4 returning skill_id,user_id,name,description,skill_status,skill_service_type;"
	err := u.db.db.QueryRow(query, newSkillName, newSkillDescription, skillId, userId).Scan(
		&updatedSkill.Id,
		&updatedSkill.UserId,
		&updatedSkill.Name,
		&updatedSkill.Description,
		&updatedSkill.Status,
	)
	if err != nil {
		return skills.Skills{}, err
	}

	return updatedSkill, nil
}

func (u *UserRepo) DeleteSkill(userId int, skillId int) (skills.Skills, error) {

	var deletedSkill skills.Skills
	query := "delete from  skills  where skill_id= $1 and user_id=$2 returning skill_id,user_id,name,description,skill_status,skill_service_type;"
	err := u.db.db.QueryRow(query, skillId, userId).Scan(
		&deletedSkill.Id,
		&deletedSkill.UserId,
		&deletedSkill.Name,
		&deletedSkill.Description,
		&deletedSkill.Status,
	)
	if err != nil {
		return skills.Skills{}, err
	}

	return deletedSkill, nil
}

func (u *UserRepo) ActivateSkill(userId int, skillId int) (skills.Skills, error) {

	var activatedSkill skills.Skills
	query := "update skills  set skill_status='active' where skill_id= $1 and user_id=$2 returning skill_id,user_id,name,description,skill_status;"
	err := u.db.db.QueryRow(query, skillId, userId).Scan(
		&activatedSkill.Id,
		&activatedSkill.UserId,
		&activatedSkill.Name,
		&activatedSkill.Description,
		&activatedSkill.Status,
	)
	if err != nil {
		return skills.Skills{}, err
	}
	return activatedSkill, nil
}

func (u *UserRepo) DectivateSkill(userId int, skillId int) (skills.Skills, error) {
	var deactivatedSkill skills.Skills
	query := "update skills  set skill_status='inactive' where skill_id= $1 and user_id=$2 returning skill_id,user_id,name,description,skill_status;"
	err := u.db.db.QueryRow(query, skillId, userId).Scan(
		&deactivatedSkill.Id,
		&deactivatedSkill.UserId,
		&deactivatedSkill.Name,
		&deactivatedSkill.Description,
		&deactivatedSkill.Status,
	)
	if err != nil {
		return skills.Skills{}, err
	}

	return deactivatedSkill, nil
}
func (u *UserRepo) CreateSession(helpToUserId int, helpFromUserId int, skillSharedId int) (helpsession.HelpSession, error) {
	var createdSession helpsession.HelpSession
	var id int
	// fmt.Println("userId", userId)
	// fmt.Println("fromUserId", fromUserId)
	query := "insert into helping_sessions(sender_id,receiver_id,skill_shared_id,time_taken,started_at,completed_at)values($1, $2, $3, $4, $5, $6) returning id;"
	fmt.Println("Skill shared id", skillSharedId)

	err := u.db.db.QueryRow(query, helpFromUserId, helpToUserId, skillSharedId, &createdSession.TimeTaken, &createdSession.StartedAt, &createdSession.CompletedAt).Scan(&id)
	if err != nil {
		return helpsession.HelpSession{}, err
	}
	createdSession.HelpToUserId = helpToUserId
	createdSession.HelpFromUserId = helpFromUserId
	createdSession.SkillSharedId = skillSharedId
	// createdSession.TimeTaken=

	createdSession.StartedAt = time.Now()
	return createdSession, nil
}

func (u *UserRepo) GetAllSessions(userId int) ([]helpsession.HelpSession, error) {

	var people []helpsession.HelpSession
	query := "select * from helping_sessions where receiver_id=$1 or sender_id=$1;"

	rows, err := u.db.db.Query(query, userId)

	if err != nil {
		fmt.Println("Error while running query        :             ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var session helpsession.HelpSession
		if err := rows.Scan(&session.Id, &session.HelpToUserId, &session.HelpFromUserId, &session.SkillSharedId, &session.TimeTaken, &session.StartedAt, &session.CompletedAt); err != nil {
			fmt.Println("Error while scanning various rows")
			return nil, err
		}
		people = append(people, session)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error before returning people[]")
		return nil, err
	}
	// fmt.Println(people)  //empty array would go in case of no users found with that skill
	return people, nil

}

func (u *UserRepo) GetSessionById(userId int, sessionId int) (helpsession.HelpSession, error) {

	var session helpsession.HelpSession
	query := "select * from helping_sessions where id=$1 and sender_id=$2 or receiver_id=$2;"
	err := u.db.db.QueryRow(query, sessionId, userId).Scan(&session.Id, &session.HelpFromUserId, &session.HelpToUserId, &session.SkillSharedId, &session.TimeTaken, &session.StartedAt, &session.CompletedAt)
	if err != nil {
		return helpsession.HelpSession{}, err
	}

	return session, nil
}
func (u *UserRepo) StopSession(userId int, sessionId int) (helpsession.HelpSession, error) {

	var stoppedSession helpsession.HelpSession
	query := "update helping_sessions set completed_at=$1 where id=$2 and (sender_id=$3 or receiver_id=$3) returning *;"
	err := u.db.db.QueryRow(query, time.Now(), sessionId, userId).Scan(&stoppedSession.Id, &stoppedSession.HelpFromUserId, &stoppedSession.HelpToUserId, &stoppedSession.SkillSharedId, &stoppedSession.TimeTaken, &stoppedSession.StartedAt, &stoppedSession.CompletedAt)
	if err != nil {
		return helpsession.HelpSession{}, err
	}

	return stoppedSession, nil
}
