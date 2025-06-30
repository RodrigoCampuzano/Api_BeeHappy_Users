package handlers

import (
	"apiusuarios/src/usuarios/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TwoFactorAuthHandler struct {
	twoFactorAuthUseCase *application.TwoFactorAuthUseCase
}

func NewTwoFactorAuthHandler(twoFactorUseCase *application.TwoFactorAuthUseCase) *TwoFactorAuthHandler {
	return &TwoFactorAuthHandler{
		twoFactorAuthUseCase: twoFactorUseCase,
	}
}

type VerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ToggleTwoFactorRequest struct {
	Estado bool `json:"estado"`
}

func (h *TwoFactorAuthHandler) ToggleTwoFactor(ctx *gin.Context) {
	var req ToggleTwoFactorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Obtener el usuario del contexto
	usuario, exists := ctx.Get("usuario")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuarioStr, ok := usuario.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	err := h.twoFactorAuthUseCase.ToggleTwoFactor(usuarioStr, req.Estado)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Configuración actualizada"})
}
