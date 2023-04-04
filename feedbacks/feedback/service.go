package feedback

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	Store(ctx context.Context, f *Feedback) (uuid.UUID, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (s *Service) Store(ctx context.Context, f *Feedback) (uuid.UUID, error) {
	f.ID = uuid.New()
	err := s.repo.Store(ctx, f)
	if err != nil {
		return uuid.Nil, err
	}
	return f.ID, nil
}
