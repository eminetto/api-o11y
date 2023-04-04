package vote

import (
	"context"
)

type Repository interface {
	Store(ctx context.Context, v *Vote) error
}
