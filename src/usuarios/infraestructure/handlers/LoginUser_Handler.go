// src/usuarios/infraestructure/handlers/LoginUser_Handler.go (ACTUALIZADO)
package handlers

import (
    "apiusuarios/src/usuarios/application"
    "net/http"
    "github.com/gin-gonic/gin"
)

type LoginUserHandler struct {
    loginUserUseCase *application.LoginUserUseCase
    tfaUseCase       *application.TFAUseCase
}

func NewLoginUserHandler(loginUseCase *application.LoginUserUseCase, tfaUseCase *application.TFAUseCase) *LoginUserHandler {
    return &LoginUserHandler{
        loginUserUseCase: loginUseCase,
        tfaUseCase:       tfaUseCase,
    }
}

type LoginRequest struct {
    Usuario    string `json:"usuario"`
    Contrasena string `json:"contrasena"`
}

type LoginTFARequest struct {
    Usuario    string `json:"usuario"`
    Contrasena string `json:"contrasena"`
    Code       string `json:"code"`
}

// Login con 2FA - Paso 1: Validar credenciales y enviar código
func (h *LoginUserHandler) LoginWithTFA(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    // Validar credenciales
    email, err := h.loginUserUseCase.ValidateCredentials(req.Usuario, req.Contrasena)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Enviar código 2FA
    err = h.tfaUseCase.SendLoginCode(*email)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando código de verificación"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Credenciales válidas. Código de verificación enviado al correo.",
        "email_masked": maskEmail(*email),
        "requires_2fa": true,
        "expires_in": 60,
    })
}

// Login con 2FA - Paso 2: Verificar código y generar token
func (h *LoginUserHandler) VerifyLoginTFA(ctx *gin.Context) {
    var req LoginTFARequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    // Validar credenciales nuevamente
    email, err := h.loginUserUseCase.ValidateCredentials(req.Usuario, req.Contrasena)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Verificar código 2FA
    err = h.tfaUseCase.VerifyLoginCode(*email, req.Code)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Generar token JWT
    token, err := h.loginUserUseCase.GenerateToken(req.Usuario)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "token": token,
        "message": "Login exitoso",
    })
}

// Login original sin 2FA (mantenido para compatibilidad)
func (h *LoginUserHandler) Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    token, err := h.loginUserUseCase.Execute(req.Usuario, req.Contrasena)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Función auxiliar para enmascarar el email
func maskEmail(email string) string {
    if len(email) < 3 {
        return "***"
    }
    
    atIndex := -1
    for i, char := range email {
        if char == '@' {
            atIndex = i
            break
        }
    }
    
    if atIndex == -1 {
        return "***"
    }
    
    if atIndex <= 2 {
        return "***" + email[atIndex:]
    }
    
    return email[:2] + "***" + email[atIndex:]
}