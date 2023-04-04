package feedback

import "context"

type Repository interface {
	Store(ctx context.Context, f *Feedback) error
}
