package vote

import "github.com/google/uuid"

type Vote struct {
	ID uuid.UUID
	Email string
	TalkName string `json:"talk_name"`
	Score int `json:"score,string"`
}

