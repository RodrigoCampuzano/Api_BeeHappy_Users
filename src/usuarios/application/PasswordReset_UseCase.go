package application

import (
	"apiusuarios/src/core/auth"
	"apiusuarios/src/usuarios/domain/entities"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetUseCase struct {
	userRepository repositories.UserRepository
	twoFactorAuth  *TwoFactorAuthUseCase
}

func NewPasswordResetUseCase(userRepo repositories.UserRepository) (*PasswordResetUseCase, error) {
	twoFactorAuth, err := NewTwoFactorAuthUseCase(userRepo)
	if err != nil {
		return nil, err
	}

	return &PasswordResetUseCase{
		userRepository: userRepo,
		twoFactorAuth:  twoFactorAuth,
	}, nil
}

func (uc *PasswordResetUseCase) RequestPasswordReset(email string) error {
	log.Printf("Iniciando solicitud de restablecimiento de contraseña para: %s", email)

	// Verificar que el usuario existe
	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error: usuario no encontrado para el email %s: %v", email, err)
		return fmt.Errorf("usuario no encontrado")
	}
	log.Printf("Usuario encontrado: %s", user.Usuario)

	// Generar y enviar código de verificación
	err = uc.twoFactorAuth.GenerateAndSendCode(email, entities.VerificationTypePasswordReset)
	if err != nil {
		log.Printf("Error al generar/enviar código de verificación: %v", err)
		return err
	}
	log.Printf("Código de verificación generado y enviado exitosamente")

	return nil
}

func (uc *PasswordResetUseCase) VerifyCodeAndResetPassword(email string, code string, newPassword string) error {
	// Verificar el código
	err := uc.twoFactorAuth.VerifyCode(email, code, entities.VerificationTypePasswordReset)
	if err != nil {
		return err
	}

	// Actualizar la contraseña
	return uc.userRepository.UpdatePassword(email, newPassword)
}

func (uc *PasswordResetUseCase) ChangePassword(email string, currentPassword string) error {
	// Verificar que el usuario existe y la contraseña actual es correcta
	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// Verificar la contraseña actual
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET no definido")
	}

	passwordHasheada := auth.HashPasswordWithSecret(currentPassword, secret)
	err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(passwordHasheada))
	if err != nil {
		return fmt.Errorf("contraseña actual incorrecta")
	}

	// Generar y enviar código de verificación
	err = uc.twoFactorAuth.GenerateAndSendCode(email, entities.VerificationTypePasswordChange)
	if err != nil {
		return err
	}

	return nil
}

func (uc *PasswordResetUseCase) VerifyCodeAndChangePassword(email string, code string, newPassword string) error {
	// Verificar el código
	err := uc.twoFactorAuth.VerifyCode(email, code, entities.VerificationTypePasswordChange)
	if err != nil {
		return err
	}

	// Actualizar la contraseña
	return uc.userRepository.UpdatePassword(email, newPassword)
}

func (uc *PasswordResetUseCase) GetUserByUsuario(usuario string) (*entities.User, error) {
	return uc.userRepository.GetUserByUsuario(usuario)
}
