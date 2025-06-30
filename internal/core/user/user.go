package user

type User struct {
	Id               int     `json:"id"`
	Username         string  `json:"username"`
	Email            string  `json:"email"`
	Password         string  `json:"password"` //omit later
	Location         string  `json:"location"`
	Availability     bool    `json:"availability"`
	AvailableCredits float64 `json:"availableCredits"`
}
type LoginRequestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` //omit later
	Id       int    `json:"id"`
}

type GetUserProfile struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Password string `json:"password"` //omit later
	Uid int `json:"uid"`
}
