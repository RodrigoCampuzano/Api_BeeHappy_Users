// src/usuarios/application/TFA_UseCase.go
package application

import (
	"apiusuarios/src/core/tfa"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
)

type TFAUseCase struct {
	userRepository repositories.UserRepository
	tfaService     *tfa.TFAService
}

func NewTFAUseCase(userRepo repositories.UserRepository) *TFAUseCase {
	return &TFAUseCase{
		userRepository: userRepo,
		tfaService:     tfa.GetTFAService(),
	}
}

func (uc *TFAUseCase) SendLoginCode(email string) error {
	// Verificar que el email existe en la base de datos
	_, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Generar y enviar código
	_, err = uc.tfaService.GenerateCode(email, "login")
	if err != nil {
		return fmt.Errorf("error generando código: %v", err)
	}

	return nil
}

func (uc *TFAUseCase) SendPasswordChangeCode(email string) error {
	// Verificar que el email existe en la base de datos
	_, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Generar y enviar código
	_, err = uc.tfaService.GenerateCode(email, "password_change")
	if err != nil {
		return fmt.Errorf("error generando código: %v", err)
	}

	return nil
}

func (uc *TFAUseCase) VerifyLoginCode(email, code string) error {
	isValid := uc.tfaService.VerifyCode(email, code, "login")
	if !isValid {
		return fmt.Errorf("código inválido o expirado")
	}
	return nil
}

func (uc *TFAUseCase) VerifyPasswordChangeCode(email, code string) error {
	isValid := uc.tfaService.VerifyCode(email, code, "password_change")
	if !isValid {
		return fmt.Errorf("código inválido o expirado")
	}
	return nil
}