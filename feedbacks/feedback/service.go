package feedback

import "github.com/google/uuid"

type UseCase interface {
	Store(f Feedback) (uuid.UUID,error)
}

type Service struct {}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Store(f Feedback) (uuid.UUID,error) {
	//@TODO create store rules, using databases or something else
	return uuid.New(),nil
}
