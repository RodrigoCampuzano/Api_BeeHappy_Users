package repositories

import (
	"apiusuarios/src/usuarios/domain/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByUsuario(usuario string) (*entities.User, error)
}