package user_test

import (
	"context"
	"github.com/eminetto/api-o11y/auth/user"
	"github.com/eminetto/api-o11y/auth/user/mocks"
	tmocks "github.com/eminetto/api-o11y/internal/telemetry/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	otel := tmocks.NewTelemetry(t)
	span := tmocks.NewSpan(t)
	span.On("RecordError", mock.Anything).Return(nil)
	span.On("SetStatus", mock.Anything, mock.Anything).Return(nil)
	span.On("End").Return(nil)
	otel.On("Start", ctx, "validatePassword").Return(ctx, span)

	s := user.NewService(repo, otel)
	u := &user.User{
		Email:    "eminetto@email.com",
		Password: "8cb2237d0679ca88db6464eac60da96345513964",
	}
	t.Run("invalid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid password", err.Error())
	})
	t.Run("valid password", func(t *testing.T) {
		err := s.ValidatePassword(ctx, u, "12345")
		assert.Nil(t, err)
	})
}
