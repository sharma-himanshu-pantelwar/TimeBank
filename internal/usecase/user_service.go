package userservice

import (
	"fmt"
	"time"
	"timebank/internal/adaptors/persistance"
	feedback "timebank/internal/core/feedback"
	helpsession "timebank/internal/core/help_session"
	"timebank/internal/core/session"
	"timebank/internal/core/skills"
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
		fmt.Println(err)
		return loginResponse, fmt.Errorf("invalid username")
	}
	loginResponse.FoundUser = foundUser

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

func (u *UserService) RegisterSkill(userId int, skill skills.Skills) (skills.Skills, error) {
	// call CreateUser function
	newSkill, err := u.userRepo.CreateSkill(userId, skill)
	return newSkill, err
}

func (u *UserService) FindPersonWithSkill(userId int, skill string) ([]user.GetUsersWithSkills, error) {
	// call CreateUser function
	foundUsers, err := u.userRepo.FindSkilledPerson(userId, skill)
	// fmt.Println(foundUsers) //empty array recieved in case of no user with that skill
	return foundUsers, err
}

func (u *UserService) RenameSkill(userId int, newSkillName string, newSkillDescription string, skillId int) (skills.Skills, error) {
	// call CreateUser function
	renamedSkill, err := u.userRepo.RenameSkill(userId, newSkillName, newSkillDescription, skillId)
	// fmt.Println(foundUsers) //empty array recieved in case of no user with that skill
	return renamedSkill, err
}

func (u *UserService) DeleteSkill(userId int, skillId int) (skills.Skills, error) {
	// call CreateUser function
	deletedSkill, err := u.userRepo.DeleteSkill(userId, skillId)
	// fmt.Println(foundUsers) //empty array recieved in case of no user with that skill
	return deletedSkill, err
}

func (u *UserService) SetActive(userId int, skillId int) (skills.Skills, error) {
	// call CreateUser function
	activatedSkill, err := u.userRepo.ActivateSkill(userId, skillId)
	// fmt.Println(foundUsers) //empty array recieved in case of no user with that skill
	return activatedSkill, err
}
func (u *UserService) SetInactive(userId int, skillId int) (skills.Skills, error) {
	// call CreateUser function
	deactivatedSkill, err := u.userRepo.DectivateSkill(userId, skillId)
	// fmt.Println(foundUsers) //empty array recieved in case of no user with that skill
	return deactivatedSkill, err
}
func (u *UserService) CreateSession(helpToUserId int, helpFromUserId int, skillSharedId int) (helpsession.HelpSession, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	createdSession, err := u.userRepo.CreateSession(helpToUserId, helpFromUserId, skillSharedId)
	if err != nil {
		return helpsession.HelpSession{}, nil
	}

	fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that

	//go routine to run and show this for corresponding go routine session reciever helper creditsOfRec

	// go routine to manage credits per second

	return createdSession, err
}

// function to manage credits per minute

func (u *UserService) GetAllSessions(helpToUserId int) ([]helpsession.HelpSession, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	allSessions, err := u.userRepo.GetAllSessions(helpToUserId)
	// fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that skill
	// fmt.Println("Error is ", err)
	return allSessions, err
}
func (u *UserService) GetSessionById(helpToUserId int, sessionId int) (helpsession.HelpSession, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	session, err := u.userRepo.GetSessionById(helpToUserId, sessionId)
	// fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that skill
	// fmt.Println("Error is ", err)
	return session, err
}
func (u *UserService) StopSession(userId int, sessionId int) (helpsession.HelpSession, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	stoppedSession, err := u.userRepo.StopSession(userId, sessionId)
	// fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that skill
	// fmt.Println("Error is ", err)
	return stoppedSession, err
}
func (u *UserService) SendFeedback(feedbackData feedback.Feedback) (feedback.Feedback, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	stoppedSession, err := u.userRepo.SendFeedback(feedbackData)
	// fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that skill
	// fmt.Println("Error is ", err)
	return stoppedSession, err
}

func (u *UserService) GetAllFeedBacks(userId int) ([]feedback.Feedback, error) {
	// call CreateUser function
	// fmt.Println("fromuserId", fromUserId)
	allFeedbacks, err := u.userRepo.GetAllFeedbacks(userId)
	// fmt.Println("created session is ", createdSession) //empty array recieved in case of no user with that skill
	// fmt.Println("Error is ", err)
	return allFeedbacks, err
}
