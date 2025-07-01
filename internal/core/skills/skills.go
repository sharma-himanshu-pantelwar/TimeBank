package skills

type Skills struct {
	Id          int    `json:"skillId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ServiceType string `json:"service_type"`
}
