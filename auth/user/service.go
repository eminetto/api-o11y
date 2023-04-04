package user

import "fmt"

type UseCase interface {
	ValidateUser(email, password string) error
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}
func (s *Service) ValidateUser(email, password string) error {
	//@TODO create validation rules, using databases or something else
	if email == "eminetto@gmail.com" && password != "1234567" {
		return fmt.Errorf("Invalid user")
	}
	return nil
}
