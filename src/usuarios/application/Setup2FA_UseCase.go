package application

import (
	"apiusuarios/src/core/auth"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
)

type Setup2FAUseCase struct {
	userRepository   repositories.UserRepository
	twoFactorService *auth.TwoFactorService
}

func NewSetup2FAUseCase(userRepo repositories.UserRepository) *Setup2FAUseCase {
	return &Setup2FAUseCase{
		userRepository:   userRepo,
		twoFactorService: auth.NewTwoFactorService(),
	}
}

// GenerateQRCode genera el código QR para configurar 2FA
func (uc *Setup2FAUseCase) GenerateQRCode(usuario string) (string, string, error) {
	// Generar nuevo secreto
	key, err := uc.twoFactorService.GenerateSecret(usuario, "API Usuarios")
	if err != nil {
		return "", "", fmt.Errorf("error generando secreto: %w", err)
	}

	// Generar código QR
	qrCode, err := uc.twoFactorService.GenerateQRCode(key)
	if err != nil {
		return "", "", fmt.Errorf("error generando código QR: %w", err)
	}

	// Guardar el secreto en la base de datos (temporalmente)
	err = uc.userRepository.UpdateTwoFactorSecret(usuario, key.Secret())
	if err != nil {
		return "", "", fmt.Errorf("error guardando secreto: %w", err)
	}

	return qrCode, key.Secret(), nil
}

// Enable2FA confirma y habilita el 2FA para el usuario
func (uc *Setup2FAUseCase) Enable2FA(usuario, token string) error {
	// Obtener el usuario para verificar el secreto
	user, err := uc.userRepository.GetUserByUsuario(usuario)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Validar el token TOTP
	if !uc.twoFactorService.ValidateToken(token, user.TwoFactorSecret) {
		return fmt.Errorf("token inválido")
	}

	// Habilitar 2FA para el usuario
	err = uc.userRepository.EnableTwoFactor(usuario)
	if err != nil {
		return fmt.Errorf("error habilitando 2FA: %w", err)
	}

	return nil
}

// Disable2FA deshabilita el 2FA para el usuario
func (uc *Setup2FAUseCase) Disable2FA(usuario, token string) error {
	// Obtener el usuario para verificar el secreto
	user, err := uc.userRepository.GetUserByUsuario(usuario)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Validar el token TOTP
	if !uc.twoFactorService.ValidateToken(token, user.TwoFactorSecret) {
		return fmt.Errorf("token inválido")
	}

	// Deshabilitar 2FA para el usuario
	err = uc.userRepository.DisableTwoFactor(usuario)
	if err != nil {
		return fmt.Errorf("error deshabilitando 2FA: %w", err)
	}

	return nil
}