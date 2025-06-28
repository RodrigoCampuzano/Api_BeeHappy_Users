package tfa

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type TFAService struct {
	codes map[string]*TFACode
}

type TFACode struct {
	Code      string
	CreatedAt time.Time
	Email     string
	Purpose   string // "login" o "password_change"
}

var tfaService *TFAService

func GetTFAService() *TFAService {
	if tfaService == nil {
		tfaService = &TFAService{
			codes: make(map[string]*TFACode),
		}
		// Limpiar códigos expirados cada 30 segundos
		go tfaService.cleanExpiredCodes()
	}
	return tfaService
}

func (t *TFAService) cleanExpiredCodes() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for email, code := range t.codes {
				if now.Sub(code.CreatedAt) > 60*time.Second {
					delete(t.codes, email)
				}
			}
		}
	}
}

func (t *TFAService) GenerateCode(email, purpose string) (string, error) {
	// Generar código de 6 dígitos
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += num.String()
	}

	// Almacenar el código
	t.codes[email] = &TFACode{
		Code:      code,
		CreatedAt: time.Now(),
		Email:     email,
		Purpose:   purpose,
	}

	// Enviar por correo
	err := t.sendCodeByEmail(email, code, purpose)
	if err != nil {
		delete(t.codes, email)
		return "", err
	}

	return code, nil
}

func (t *TFAService) VerifyCode(email, code, purpose string) bool {
	storedCode, exists := t.codes[email]
	if !exists {
		return false
	}

	// Verificar si ha expirado (60 segundos)
	if time.Since(storedCode.CreatedAt) > 60*time.Second {
		delete(t.codes, email)
		return false
	}

	// Verificar si el código y propósito coinciden
	if storedCode.Code == code && storedCode.Purpose == purpose {
		delete(t.codes, email) // Eliminar después del uso
		return true
	}

	return false
}

func (t *TFAService) sendCodeByEmail(email, code, purpose string) error {
	godotenv.Load()
	
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("configuración SMTP incompleta")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	var subject, body string
	switch purpose {
	case "login":
		subject = "Código de verificación para inicio de sesión"
		body = fmt.Sprintf(`
Estimado usuario,

Su código de verificación para iniciar sesión es: %s

Este código expira en 60 segundos.

Si no solicitó este código, ignore este mensaje.

Saludos,
Equipo de Seguridad
`, code)
	case "password_change":
		subject = "Código de verificación para cambio de contraseña"
		body = fmt.Sprintf(`
Estimado usuario,

Su código de verificación para cambiar su contraseña es: %s

Este código expira en 60 segundos.

Si no solicitó este código, ignore este mensaje y considere cambiar su contraseña.

Saludos,
Equipo de Seguridad
`, code)
	default:
		return fmt.Errorf("propósito no válido")
	}

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", email, subject, body))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, msg)
}