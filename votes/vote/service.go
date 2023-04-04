package vote

import "github.com/google/uuid"

type UseCase interface {
	Store(v Vote) (uuid.UUID,error)
}

type Service struct {}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Store(v Vote) (uuid.UUID,error) {
	//@TODO create store rules, using databases or something else
	return uuid.New(),nil
}
