package user_test

import (
	"auth/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateUser(t *testing.T){
	service := user.NewService()
	t.Run("invalid user", func(t *testing.T) {
		err := service.ValidateUser("eminetto@gmail.com", "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "Invalid user", err.Error())
	})
	t.Run("valid user", func(t *testing.T) {
		err := service.ValidateUser("eminetto@gmail.com", "1234567")
		assert.Nil(t, err)
	})
}
