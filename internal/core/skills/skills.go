package skills

type Skills struct {
	Id          int    `json:"skillId"`
	Name        string `json:"name"`
	UserId      int    `json:"userId"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ServiceType string `json:"service_type"`
}
