package user

import (
	"context"
	"crypto/sha1"
	"fmt"
)

type UseCase interface {
	ValidateUser(ctx context.Context, email, password string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) ValidateUser(ctx context.Context, email, password string) error {
	u, err := s.repo.Get(ctx, email)
	if err != nil {
		return err
	}
	return validatePassword(u, password)
}

func validatePassword(u *User, password string) error {
	h := sha1.New()
	h.Write([]byte(password))
	p := fmt.Sprintf("%x", h.Sum(nil))
	if p != u.Password {
		return fmt.Errorf("invalid user")
	}
	return nil
}
