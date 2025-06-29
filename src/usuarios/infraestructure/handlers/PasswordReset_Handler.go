package handlers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	passwordResetUseCase *application.PasswordResetUseCase
}

func NewPasswordResetHandler(passwordResetUseCase *application.PasswordResetUseCase) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetUseCase: passwordResetUseCase,
	}
}

func (h *PasswordResetHandler) RequestReset(ctx *gin.Context) {
	var req entities.SwaggerPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.passwordResetUseCase.RequestPasswordReset(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities.SwaggerPasswordResetResponse{
		Message: "Código de verificación enviado",
	})
}

func (h *PasswordResetHandler) ResetPassword(ctx *gin.Context) {
	var req entities.SwaggerResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.passwordResetUseCase.VerifyCodeAndResetPassword(req.Email, req.Code, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities.SwaggerResetPasswordResponse{
		Message: "Contraseña actualizada",
	})
}

type RequestResetRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *PasswordResetHandler) RequestChangePassword(ctx *gin.Context) {
	var req entities.SwaggerChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Obtener el usuario del token JWT
	usuario := ctx.GetString("usuario")
	if usuario == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el email del usuario
	var user *entities.User
	var err error
	user, err = h.passwordResetUseCase.GetUserByUsuario(usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	email := user.Correo_electronico

	err = h.passwordResetUseCase.ChangePassword(email, req.CurrentPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities.SwaggerChangePasswordResponse{
		Message: "Código de verificación enviado",
	})
}

func (h *PasswordResetHandler) ChangePassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.passwordResetUseCase.VerifyCodeAndChangePassword(req.Email, req.Code, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada"})
}
