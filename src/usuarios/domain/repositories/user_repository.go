package repositories

import (
	"apiusuarios/src/usuarios/domain/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByUsuario(usuario string) (*entities.User, error)
	UpdateTwoFactorSecret(usuario, secret string) error
	EnableTwoFactor(usuario string) error
	DisableTwoFactor(usuario string) error
}