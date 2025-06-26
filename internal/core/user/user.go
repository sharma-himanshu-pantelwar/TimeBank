package user

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` //omit later
	Uid      int    `json:"uid"`
}
type LoginRequestUser struct {
	Username string `json:"username"`
	Password string `json:"password"` //omit later
	Uid      int    `json:"uid"`
}

type GetUserProfile struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Password string `json:"password"` //omit later
	Uid int `json:"uid"`
}
