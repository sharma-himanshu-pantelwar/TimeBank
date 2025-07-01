package helpsession

import "time"

type HelpSession struct {
	ToUser        int       `json:"helpNeededBy` //this person is getting helped
	FromUser      int       `json:"helpGivenBy`  //this person is  helping
	SkillSharedId int       `json:"skillsharedid"`
	TimeTaken     float64   `json:"timeTaken"`
	StartedAt     time.Time `json:"sessionStartedAt"`
	CompletedAt   time.Time `json:"sessionCompletedAt"`
}
