package application

import (
	"apiusuarios/src/usuarios/domain/entities"
	"apiusuarios/src/usuarios/domain/repositories"
)

type CreateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewCreateUserUseCase(userRepo repositories.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepo,
	}
}

func (uc *CreateUserUseCase) Execute(data *entities.User) error {
	return uc.userRepository.CreateUser(data)
}