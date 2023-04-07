package feedback

import (
	"context"
	"github.com/eminetto/api-o11y/internal/telemetry"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

type UseCase interface {
	Store(ctx context.Context, f *Feedback) (uuid.UUID, error)
}

type Service struct {
	repo      Repository
	telemetry telemetry.Telemetry
}

func NewService(repo Repository, otel telemetry.Telemetry) *Service {
	return &Service{
		repo:      repo,
		telemetry: otel,
	}
}
func (s *Service) Store(ctx context.Context, f *Feedback) (uuid.UUID, error) {
	ctx, span := s.telemetry.Start(ctx, "service")
	defer span.End()
	f.ID = uuid.New()
	err := s.repo.Store(ctx, f)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, err
	}
	return f.ID, nil
}
