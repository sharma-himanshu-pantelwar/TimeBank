package feedback

import "time"

type Feedback struct {
	Id        int       `json:"id"`
	SessionId int       `json:"session_id"`
	RaterId   int       `json:"rater_id"`
	RateeId   int       `json:"ratee_id"`
	Rating    int       `json:"rating"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}
