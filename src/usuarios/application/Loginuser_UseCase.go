package application

import (
	"apiusuarios/src/core/auth"
	"apiusuarios/src/usuarios/domain/repositories"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepository   repositories.UserRepository
	twoFactorService *auth.TwoFactorService
}

func NewLoginUserUseCase(userRepo repositories.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository:   userRepo,
		twoFactorService: auth.NewTwoFactorService(),
	}
}

func (uc *LoginUserUseCase) Execute(usuario, contrasena, tokenTOTP string) (string, error) {
	user, err := uc.userRepository.GetUserByUsuario(usuario)
	if err != nil {
		return "", fmt.Errorf("usuario o contraseña incorrectos")
	}

	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET no definido")
	}

	// Verificar contraseña
	passwordHasheada := auth.HashPasswordWithSecret(contrasena, secret)
	err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(passwordHasheada))
	if err != nil {
		return "", fmt.Errorf("usuario o contraseña incorrectos")
	}

	// Verificar 2FA si está habilitado
	if user.TwoFactorEnabled {
		if tokenTOTP == "" {
			return "", fmt.Errorf("token 2FA requerido")
		}

		if !uc.twoFactorService.ValidateToken(tokenTOTP, user.TwoFactorSecret) {
			return "", fmt.Errorf("token 2FA inválido")
		}
	}

	// Generar JWT
	claims := jwt.MapClaims{
		"usuario": user.Usuario,
		"rol":     user.Rol,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("no se pudo generar el token")
	}

	return tokenString, nil
}