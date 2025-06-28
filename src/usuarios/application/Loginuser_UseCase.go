// src/usuarios/application/Loginuser_UseCase.go (ACTUALIZADO)
package application

import (
	"apiusuarios/src/core/auth"
    "apiusuarios/src/usuarios/domain/repositories"
    "golang.org/x/crypto/bcrypt"
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/joho/godotenv"
    "fmt"
)

type LoginUserUseCase struct {
    userRepository repositories.UserRepository
}

func NewLoginUserUseCase(userRepo repositories.UserRepository) *LoginUserUseCase {
    return &LoginUserUseCase{
        userRepository: userRepo,
    }
}

// Validar credenciales sin generar token (para proceso de 2FA)
func (uc *LoginUserUseCase) ValidateCredentials(usuario, contrasena string) (*string, error) {
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

    // Retornar el email para el proceso de 2FA
    return &user.Correo_electronico, nil
}

// Generar token después de validación exitosa de 2FA
func (uc *LoginUserUseCase) GenerateToken(usuario string) (string, error) {
    user, err := uc.userRepository.GetUserByUsuario(usuario)
    if err != nil {
        return "", fmt.Errorf("usuario no encontrado")
    }

    godotenv.Load()
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET no definido")
    }

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

// Método original mantenido para compatibilidad (sin 2FA)
func (uc *LoginUserUseCase) Execute(usuario, contrasena string) (string, error) {
    user, err := uc.userRepository.GetUserByUsuario(usuario)
    if err != nil {
        return "", fmt.Errorf("usuario o contraseña incorrectos")
    }

    godotenv.Load()
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET no definido")
    }

	passwordHasheada := auth.HashPasswordWithSecret(contrasena, secret)
    err = bcrypt.CompareHashAndPassword([]byte(user.Contrasena), []byte(passwordHasheada))
    if err != nil {
        return "", fmt.Errorf("usuario o contraseña incorrectos")
    }

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