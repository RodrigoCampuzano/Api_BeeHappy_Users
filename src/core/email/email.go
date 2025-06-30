package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type EmailService struct {
	from     string
	password string
	host     string
	port     string
}

func NewEmailService() (*EmailService, error) {
	log.Println("Inicializando servicio de correo")
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error cargando variables de entorno: %v", err)
		return nil, fmt.Errorf("error cargando variables de entorno: %v", err)
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")

	// Verificar que todas las variables necesarias estén configuradas
	if from == "" || password == "" || host == "" || port == "" {
		log.Printf("Faltan variables de entorno. FROM: %v, HOST: %v, PORT: %v, PASSWORD: %v",
			from != "", host != "", port != "", password != "")
		return nil, fmt.Errorf("faltan variables de entorno para el servicio de correo")
	}

	log.Printf("Servicio de correo configurado con FROM: %s, HOST: %s, PORT: %s", from, host, port)
	return &EmailService{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}, nil
}

func (s *EmailService) SendVerificationCode(to string, code string, purpose string) error {
	log.Printf("Preparando envío de código a %s (propósito: %s)", to, purpose)

	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var subject string
	var body string

	switch purpose {
	case "login":
		subject = "Código de verificación para inicio de sesión"
		body = fmt.Sprintf("Tu código de verificación para iniciar sesión es: %s\nEste código expirará en 2 minutos.", code)
	case "reset":
		subject = "Código de verificación para restablecer contraseña"
		body = fmt.Sprintf("Tu código de verificación para restablecer tu contraseña es: %s\nEste código expirará en 2 minutos.", code)
	case "change":
		subject = "Código de verificación para cambiar contraseña"
		body = fmt.Sprintf("Tu código de verificación para cambiar tu contraseña es: %s\nEste código expirará en 2 minutos.", code)
	}

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body)

	log.Printf("Intentando enviar correo a través de %s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Error enviando correo: %v", err)
		return err
	}
	log.Printf("Correo enviado exitosamente a %s", to)

	return nil
}
