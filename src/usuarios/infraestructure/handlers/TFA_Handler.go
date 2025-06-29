// src/usuarios/infraestructure/handlers/TFA_Handler.go
package handlers

import (
	"apiusuarios/src/usuarios/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TFAHandler struct {
	tfaUseCase            *application.TFAUseCase
	changePasswordUseCase *application.ChangePasswordUseCase
}

func NewTFAHandler(tfaUseCase *application.TFAUseCase, changePasswordUseCase *application.ChangePasswordUseCase) *TFAHandler {
	return &TFAHandler{
		tfaUseCase:            tfaUseCase,
		changePasswordUseCase: changePasswordUseCase,
	}
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// Enviar código para login
func (h *TFAHandler) SendLoginCode(ctx *gin.Context) {
	var req SendCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email requerido y válido"})
		return
	}

	err := h.tfaUseCase.SendLoginCode(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Código enviado al correo electrónico",
		"expires_in": 60,
	})
}

// Verificar código para login
func (h *TFAHandler) VerifyLoginCode(ctx *gin.Context) {
	var req VerifyCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email y código de 6 dígitos requeridos"})
		return
	}

	err := h.tfaUseCase.VerifyLoginCode(req.Email, req.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Código verificado correctamente",
		"verified": true,
	})
}

// Enviar código para cambio de contraseña
func (h *TFAHandler) SendPasswordChangeCode(ctx *gin.Context) {
	var req SendCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email requerido y válido"})
		return
	}

	err := h.tfaUseCase.SendPasswordChangeCode(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Código enviado al correo electrónico",
		"expires_in": 60,
	})
}

// Cambiar contraseña con verificación de código
func (h *TFAHandler) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email, código de 6 dígitos y nueva contraseña requeridos"})
		return
	}

	// Primero verificar el código
	err := h.tfaUseCase.VerifyPasswordChangeCode(req.Email, req.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Si el código es válido, cambiar la contraseña
	err = h.changePasswordUseCase.Execute(req.Email, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Contraseña cambiada exitosamente",
	})
}