package helpsession

import "time"

type HelpSession struct {
	HelpToUserId   int       `json:"helpNeededBy"` //this person is getting helped
	HelpFromUserId int       `json:"helpGivenBy"`  //this person is  helping
	SkillSharedId  int       `json:"skillSharedId"`
	TimeTaken      float64   `json:"timeTaken"`
	StartedAt      time.Time `json:"sessionStartedAt"`
	CompletedAt    time.Time `json:"sessionCompletedAt"`
}
