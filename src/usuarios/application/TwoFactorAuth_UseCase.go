package application

import (
	"apiusuarios/src/core/email"
	"apiusuarios/src/usuarios/domain/entities"
	"apiusuarios/src/usuarios/domain/repositories"
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type TwoFactorAuthUseCase struct {
	userRepository repositories.UserRepository
	emailService   *email.EmailService
}

func NewTwoFactorAuthUseCase(userRepo repositories.UserRepository) (*TwoFactorAuthUseCase, error) {
	emailService, err := email.NewEmailService()
	if err != nil {
		return nil, err
	}

	return &TwoFactorAuthUseCase{
		userRepository: userRepo,
		emailService:   emailService,
	}, nil
}

func (uc *TwoFactorAuthUseCase) GenerateAndSendCode(email string, tipo string) error {
	log.Printf("Generando código de verificación para %s (tipo: %s)", email, tipo)

	// Generar código de 6 dígitos que no empiece con cero
	code := make([]byte, 5)
	for {
		_, err := rand.Read(code)
		if err != nil {
			log.Printf("Error generando bytes aleatorios: %v", err)
			return err
		}

		// Asegurar que el primer dígito no sea cero
		firstDigit := (int(code[0]) % 9) + 1 // 1-9
		otherDigits := int(code[1])<<24 | int(code[2])<<16 | int(code[3])<<8 | int(code[4])
		num := firstDigit*100000 + (otherDigits % 100000)

		if num >= 100000 && num <= 999999 {
			verificationCode := fmt.Sprintf("%d", num)
			log.Printf("Código generado: %s", verificationCode)

			// Crear registro de código de verificación
			verificationRecord := &entities.VerificationCode{
				CorreoElectronico: email,
				Codigo:            verificationCode,
				FechaExpiracion:   time.Now().Add(2 * time.Minute).Format(time.RFC3339),
				Tipo:              tipo,
			}

			// Guardar código en la base de datos
			err = uc.userRepository.SaveVerificationCode(verificationRecord)
			if err != nil {
				log.Printf("Error guardando código en la base de datos: %v", err)
				return err
			}
			log.Printf("Código guardado en la base de datos")

			// Enviar código por correo
			err = uc.emailService.SendVerificationCode(email, verificationCode, tipo)
			if err != nil {
				log.Printf("Error enviando correo: %v", err)
				return err
			}
			log.Printf("Correo enviado exitosamente")

			return nil
		}
	}
}

func (uc *TwoFactorAuthUseCase) VerifyCode(email string, code string, tipo string) error {
	storedCode, err := uc.userRepository.GetVerificationCode(email, tipo)
	if err != nil {
		return fmt.Errorf("código inválido o expirado")
	}

	if storedCode.Codigo != code {
		return fmt.Errorf("código incorrecto")
	}

	// Si el código es válido, lo eliminamos para que no pueda ser reutilizado
	return uc.userRepository.DeleteVerificationCode(storedCode.ID)
}

func (uc *TwoFactorAuthUseCase) ToggleTwoFactor(usuario string, estado bool) error {
	// Obtener el usuario por nombre de usuario
	user, err := uc.userRepository.GetUserByUsuario(usuario)
	if err != nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Actualizar el estado de verificación en dos pasos
	return uc.userRepository.UpdateVerificacionDospasos(user.ID, estado)
}
