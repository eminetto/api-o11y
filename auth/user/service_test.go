package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	u := &User{
		Email:    "eminetto@email.com",
		Password: "8cb2237d0679ca88db6464eac60da96345513964",
	}
	t.Run("invalid password", func(t *testing.T) {
		err := validatePassword(u, "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid user", err.Error())
	})
	t.Run("valid password", func(t *testing.T) {
		err := validatePassword(u, "12345")
		assert.Nil(t, err)
	})
}
