package userservice

import (
	"fmt"
	"time"
	"timebank/internal/adaptors/persistance"
	"timebank/internal/core/session"
	"timebank/internal/core/user"
	"timebank/pkg/generatejwt"
	hashpassword "timebank/pkg/hashPassword"

	"golang.org/x/crypto/bcrypt"
)

// ^ Write function for user registration

// 1. Create struct for UserService which refers to informai
type UserService struct {
	userRepo    persistance.UserRepo //reference to user repo which is basically a reference to a database
	sessionRepo persistance.SessionRepo
}

func NewUserService(userRepo persistance.UserRepo, sessionRepo persistance.SessionRepo) UserService {
	return UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

// Register user method
// method is applied over u(ref. to Userservice ref. to userRepo which is basically ref. to DB  )
func (u *UserService) RegisterUser(user user.User) (user.User, error) {
	// call CreateUser function
	newUser, err := u.userRepo.CreateUser(user)
	return newUser, err
}

type LoginResponse struct {
	FoundUser   user.User       //user found with their struct
	TokenString string          //token value string
	TokenExpire time.Time       //expiry date of token
	Session     session.Session //session found with id
}

// this function takes ref. to UserService(ref. to databases for userRepo and SessionRepo) as parameter and returns LoginResponse,error(if any)
func (u *UserService) LoginUser(requestUser user.LoginRequestUser) (LoginResponse, error) {
	// fmt.Println(requestUser)//getting req body perfectly
	loginResponse := LoginResponse{}
	// fmt.Println("Would be searching for ", requestUser.Username)
	foundUser, err := u.userRepo.GetUser(requestUser.Username) //this will find user based on their username
	// fmt.Println("Found user :               ", foundUser) returns username pass uid
	if err != nil {
		return loginResponse, fmt.Errorf("invalid username")
	}
	loginResponse.FoundUser = foundUser

	// fmt.Println("Login response ", loginResponse)
	// fmt.Println("loginResposeFound user ", loginResponse.FoundUser)
	// fmt.Println("Found user.Uid   ", foundUser.Uid)

	// fmt.Println("Request user   ", requestUser)
	// fmt.Println("Request pass   ", requestUser.Password)

	// fmt.Println("Found user is : ", foundUser)

	if err := matchPassword(foundUser, requestUser.Password); err != nil {
		return loginResponse, fmt.Errorf("invalid username or password")
	}

	tokenString, tokenExpire, err := generatejwt.GenerateJWT(foundUser.Id)
	loginResponse.TokenString = tokenString
	loginResponse.TokenExpire = tokenExpire

	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate jwt token")
	}

	session, err := generatejwt.GenerateSession(foundUser.Id)
	loginResponse.Session = session
	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate session")
	}

	// create session
	err = u.sessionRepo.CreateSession(session)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to create session")
	}
	// fmt.Println(loginResponse)
	return loginResponse, nil

}

func matchPassword(user user.User, password string) error {
	// fmt.Println("Actual password for user is ", user.Password)
	// fmt.Println("Recieved Password for  user is ", password)

	err := hashpassword.CheckPassword(user.Password, password)
	if err != nil {
		return fmt.Errorf("unable to match password")
	}
	return nil
}

func (u *UserService) GetUserById(id int) (user.GetUserProfile, error) {
	newUser, err := u.userRepo.GetUserById(id)
	return newUser, err
}

func (u *UserService) LogoutUser(id int) error {
	err := u.sessionRepo.DeleteSession(id)
	return err
}

func (u *UserService) GetJWTFromSession(sess string) (string, time.Time, error) {
	var tokenString string
	var tokenExpire time.Time
	session, err := u.sessionRepo.GetSession(sess)
	if err != nil {
		return tokenString, tokenExpire, err
	}
	err = matchSessionToken(sess, session.TokenHash)
	if err != nil {
		return tokenString, tokenExpire, err
	}
	tokenString, tokenExpire, err = generatejwt.GenerateJWT(session.Uid)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	return tokenString, tokenExpire, nil
}

func matchSessionToken(id string, tokenHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(id))
	if err != nil {
		fmt.Println(err, "unable to match session token")
		return err
	}
	return nil
}
