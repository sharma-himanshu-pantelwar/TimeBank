package skills

type Skills struct {
	Id              int     `json:"skillId"`
	Name            string  `json:"name"`
	UserId          int     `json:"userId"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	MinTimeRequired float64 `json:"minTimeRequired"`
}
