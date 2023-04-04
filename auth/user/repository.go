package user

import "context"

type Repository interface {
	Get(ctx context.Context, email string) (*User, error)
}
