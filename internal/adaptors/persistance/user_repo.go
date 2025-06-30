package persistance

import (
	"fmt"
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
