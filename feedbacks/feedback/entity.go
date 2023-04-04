package feedback

import "github.com/google/uuid"

type Feedback struct {
	ID uuid.UUID
	Email string
	Title string `json:"title"`
	Body string `json:"body"`
}
