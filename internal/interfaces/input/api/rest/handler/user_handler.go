package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	feedback "timebank/internal/core/feedback"
	helpsession "timebank/internal/core/help_session"
	"timebank/internal/core/skills"
	"timebank/internal/core/user"
	userservice "timebank/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService userservice.UserService //this service returns newUser,err newuser ->User struct
}

func NewUserHandler(usecase userservice.UserService) UserHandler {
	return UserHandler{
		userService: usecase,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user user.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	registeredUser, err := u.userService.RegisterUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user = registeredUser
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// var requestUser user.User
	var loginRequestUser user.LoginRequestUser

	if err := json.NewDecoder(r.Body).Decode(&loginRequestUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// fmt.Println("After decode    :           ", loginRequestUser)
	loginResponse, err := u.userService.LoginUser(loginRequestUser)
	// fmt.Println(loginResponse) //getting username,hashedPassword,uid in loginResponse
	if err != nil {
		// fmt.Println("Got an error while recieving loginResponse\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// cookie
	atCookie := http.Cookie{
		Name:     "at",
		Value:    loginResponse.TokenString,
		Expires:  loginResponse.TokenExpire,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    loginResponse.Session.Id.String(),
		Expires:  loginResponse.Session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &sessCookie)

	// fmt.Println("Usernamre in loginResponse is :", loginResponse.FoundUser.Username)//not coming upto here
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", loginResponse.FoundUser.Username)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"message": "Successful login"})
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	returnedUser, err := u.userService.GetUserById(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", returnedUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(returnedUser)

}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	err := u.userService.LogoutUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	atCookie := http.Cookie{
		Name:     "at",
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)
	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &sessCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "successful logout"})

}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, expireTime, err := u.userService.GetJWTFromSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	atCookie := http.Cookie{
		Name:     "at",
		Value:    tokenString,
		Expires:  expireTime,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "cookie refreshed succesfully"})
}

// Register Skills
func (u *UserHandler) AddSkills(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	var newSkills skills.Skills

	fmt.Println("user id while registration is", userId)
	if err := json.NewDecoder(r.Body).Decode(&newSkills); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	registeredSkill, err := u.userService.RegisterSkill(userId, newSkills)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	skill := registeredSkill
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skill)
}

func (u *UserHandler) FindSkilledPerson(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	skill := chi.URLParam(r, "skill")
	// 	var newSkills skills.Skills

	// 	if err := json.NewDecoder(r.Body).Decode(&newSkills); err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}
	foundUsersWithSkill, err := u.userService.FindPersonWithSkill(userId, skill)
	if err != nil {
		fmt.Println("Error after calling FindPersonWithSkill from handler")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := foundUsersWithSkill
	if len(foundUsersWithSkill) == 0 {
		response = []user.GetUsersWithSkills{}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (u *UserHandler) RenameSkill(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	skillIdStr := chi.URLParam(r, "skillId")
	skillId, err := strconv.Atoi(skillIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}
	// 	var newSkills skills.Skills
	type RenameRequest struct {
		NewNameForSkill     string `json:"newName"`
		NewDescriptionSkill string `json:"newDescription"`
	}

	var alteredSkill RenameRequest
	if err := json.NewDecoder(r.Body).Decode(&alteredSkill); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	renamedSkillResponse, err := u.userService.RenameSkill(userId, alteredSkill.NewNameForSkill, alteredSkill.NewDescriptionSkill, skillId)
	if err != nil {
		fmt.Println("Error after calling FindPersonWithSkill from handler")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(renamedSkillResponse)
}

func (u *UserHandler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	skillIdStr := chi.URLParam(r, "skillId")
	skillId, err := strconv.Atoi(skillIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}
	// 	var newSkills skills.Skills

	// var deletedSkill skills.Skills
	// if err := json.NewDecoder(r.Body).Decode(&deletedSkill); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	deletedSkillResponse, err := u.userService.DeleteSkill(userId, skillId)
	if err != nil {
		fmt.Println("Error after calling FindPersonWithSkill from handler")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedSkillResponse)
}

func (u *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	skillIdStr := chi.URLParam(r, "skillId")
	skillId, err := strconv.Atoi(skillIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}
	activateSkillResponse, err := u.userService.SetActive(userId, skillId)
	if err != nil {
		fmt.Println("Error after calling FindPersonWithSkill from handler")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activateSkillResponse)
}

func (u *UserHandler) SetInactive(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	skillIdStr := chi.URLParam(r, "skillId")
	skillId, err := strconv.Atoi(skillIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	deactivateSkillResponse, err := u.userService.SetInactive(userId, skillId)
	if err != nil {
		fmt.Println("Error after calling FindPersonWithSkill from handler")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deactivateSkillResponse)
}

func (u *UserHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var sessionData helpsession.HelpSession
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionData.HelpToUserId = userId
	// userId-> logged in user id

	// fmt.Println("fromUserId", fromUserId)
	// fmt.Printf("Sending user %v of type %T", userId, userId)

	//fmt.Println("from user", sessionData.FromUser)
	if sessionData.HelpFromUserId == sessionData.HelpToUserId {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("user can not take help from himself")
		return
	}

	deactivateSkillResponse, err := u.userService.CreateSession(sessionData.HelpToUserId, sessionData.HelpFromUserId, sessionData.SkillSharedId)

	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deactivateSkillResponse)
}

func (u *UserHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	var sessionData helpsession.HelpSession
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionData.HelpToUserId = userId
	// userId-> logged in user id

	// fmt.Println("fromUserId", fromUserId)
	// fmt.Printf("Sending user %v of type %T", userId, userId)

	//fmt.Println("from user", sessionData.FromUser)
	allSessionsResponse, err := u.userService.GetAllSessions(sessionData.HelpToUserId)

	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allSessionsResponse)
}

func (u *UserHandler) GetSessionById(w http.ResponseWriter, r *http.Request) {
	var sessionData helpsession.HelpSession
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	sessionId, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	sessionsResponse, err := u.userService.GetSessionById(userId, sessionId)

	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionsResponse)
}
func (u *UserHandler) StopSession(w http.ResponseWriter, r *http.Request) {
	var sessionData helpsession.HelpSession
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&sessionData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sessionIdStr := chi.URLParam(r, "sessionId")
	sessionId, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	stoppedSessionResponse, err := u.userService.StopSession(userId, sessionId)

	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := foundUsersWithSkill
	// if len(foundUsersWithSkill) == 0 {
	// 	response = []user.GetUsersWithSkills{}
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stoppedSessionResponse)
}

func (u *UserHandler) GiveFeedback(w http.ResponseWriter, r *http.Request) {
	var feedbackData feedback.Feedback
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&feedbackData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	rateeIdStr := chi.URLParam(r, "feedbackTo")
	rateeId, err := strconv.Atoi(rateeIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	feedbackData.RateeId = rateeId
	feedbackData.RaterId = userId
	feedbackDataResponse, err := u.userService.SendFeedback(feedbackData)
	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbackDataResponse)

}
func (u *UserHandler) GetFeedbacksForMe(w http.ResponseWriter, r *http.Request) {
	// var feedbackData feedback.Feedback
	userId, ok := r.Context().Value("user").(int)
	// fmt.Println("user id is ", userId)//2
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	// if err := json.NewDecoder(r.Body).Decode(&feedbackData); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	feedbackDataResponse, err := u.userService.GetAllFeedBacks(userId)
	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := feedbackDataResponse
	if len(feedbackDataResponse) == 0 {
		response = []feedback.Feedback{}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func (u *UserHandler) GetFeedbacksForUser(w http.ResponseWriter, r *http.Request) {
	// var feedbackData feedback.Feedback
	_, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	feedbackForUserIdStr := chi.URLParam(r, "forUser")
	feedbackForUserId, err := strconv.Atoi(feedbackForUserIdStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	feedbackDataResponse, err := u.userService.GetAllFeedBacks(feedbackForUserId)
	if err != nil {
		// fmt.Println("Error after CreateSession from handler", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := feedbackDataResponse
	if len(feedbackDataResponse) == 0 {
		response = []feedback.Feedback{}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
