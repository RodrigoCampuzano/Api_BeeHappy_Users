package repositories

import (
	"apiusuarios/src/usuarios/domain/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
}