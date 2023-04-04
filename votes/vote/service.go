package vote

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	Store(ctx context.Context, v *Vote) (uuid.UUID, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) Store(ctx context.Context, v *Vote) (uuid.UUID, error) {
	v.ID = uuid.New()
	err := s.repo.Store(ctx, v)
	if err != nil {
		return uuid.Nil, err
	}
	return v.ID, nil
}
