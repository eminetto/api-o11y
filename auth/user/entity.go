package user

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
}
