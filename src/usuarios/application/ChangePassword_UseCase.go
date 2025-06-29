// src/usuarios/application/ChangePassword_UseCase.go
package application

import (
	"apiusuarios/src/core/auth"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordUseCase struct {
	userRepository repositories.UserRepository
}

func NewChangePasswordUseCase(userRepo repositories.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepository: userRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(email, newPassword string) error {
	// Cargar variables de entorno
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET no definido")
	}

	// Verificar que el usuario existe
	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Hashear la nueva contraseña
	passwordHasheada := auth.HashPasswordWithSecret(newPassword, secret)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordHasheada), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error encriptando la contraseña: %v", err)
	}

	// Actualizar la contraseña en la base de datos
	err = uc.userRepository.UpdatePassword(user.Usuario, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error actualizando la contraseña: %v", err)
	}

	return nil
}