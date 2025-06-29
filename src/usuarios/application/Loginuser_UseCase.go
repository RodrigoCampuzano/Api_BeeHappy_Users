package application

import (
	"apiusuarios/src/core/auth"
	"apiusuarios/src/usuarios/domain/entities"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepository repositories.UserRepository
	twoFactorAuth  *TwoFactorAuthUseCase
}

func NewLoginUserUseCase(userRepo repositories.UserRepository) (*LoginUserUseCase, error) {
	twoFactorAuth, err := NewTwoFactorAuthUseCase(userRepo)
	if err != nil {
		return nil, err
	}

	return &LoginUserUseCase{
		userRepository: userRepo,
		twoFactorAuth:  twoFactorAuth,
	}, nil
}

type LoginResult struct {
	RequireTwoFactor bool   `json:"require_two_factor"`
	Token            string `json:"token,omitempty"`
	Email            string `json:"email,omitempty"`
	Message          string `json:"message,omitempty"`
}

func (uc *LoginUserUseCase) Execute(usuario, contrasena string) (*LoginResult, error) {
	user, err := uc.userRepository.GetUserByUsuario(usuario)
	if err != nil {
		return nil, fmt.Errorf("usuario o contraseña incorrectos")
	}

	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET no definido")
	}

	passwordHasheada := auth.HashPasswordWithSecret(contrasena, secret)
	err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(passwordHasheada))
	if err != nil {
		return nil, fmt.Errorf("usuario o contraseña incorrectos")
	}

	// Si el usuario tiene activada la verificación en dos pasos
	if user.VerificacionDospasos {
		// Generar y enviar código
		err = uc.twoFactorAuth.GenerateAndSendCode(user.Correo_electronico, entities.VerificationTypeLogin)
		if err != nil {
			return nil, err
		}

		return &LoginResult{
			RequireTwoFactor: true,
			Email:            user.Correo_electronico,
			Message:          "Se ha enviado un código de verificación a su correo electrónico",
		}, nil
	}

	// Si no requiere verificación en dos pasos, generamos el token directamente
	claims := jwt.MapClaims{
		"usuario": user.Usuario,
		"rol":     user.Rol,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("no se pudo generar el token")
	}

	return &LoginResult{
		RequireTwoFactor: false,
		Token:            tokenString,
		Message:          "Inicio de sesión exitoso",
	}, nil
}

func (uc *LoginUserUseCase) VerifyTwoFactorAndLogin(email string, code string) (*LoginResult, error) {
	// Verificar el código
	err := uc.twoFactorAuth.VerifyCode(email, code, entities.VerificationTypeLogin)
	if err != nil {
		return nil, err
	}

	// Obtener el usuario
	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Verificar que el usuario tenga 2FA activado
	if !user.VerificacionDospasos {
		return nil, fmt.Errorf("la verificación en dos pasos no está activada para este usuario")
	}

	// Generar el token
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET no definido")
	}

	claims := jwt.MapClaims{
		"usuario": user.Usuario,
		"rol":     user.Rol,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("no se pudo generar el token")
	}

	return &LoginResult{
		RequireTwoFactor: false,
		Token:            tokenString,
		Message:          "Código verificado correctamente. Sesión iniciada.",
	}, nil
}
