package repositories

import (
	"apiusuarios/src/usuarios/domain/entities"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByUsuario(usuario string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdatePassword(email string, newPassword string) error
	UpdateVerificacionDospasos(userID int, estado bool) error
	SaveVerificationCode(code *entities.VerificationCode) error
	GetVerificationCode(email string, tipo string) (*entities.VerificationCode, error)
	DeleteVerificationCode(id int) error
}
