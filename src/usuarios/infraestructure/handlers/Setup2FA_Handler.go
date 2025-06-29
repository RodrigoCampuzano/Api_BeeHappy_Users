package handlers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Setup2FAHandler struct {
	setup2FAUseCase *application.Setup2FAUseCase
}

func NewSetup2FAHandler(useCase *application.Setup2FAUseCase) *Setup2FAHandler {
	return &Setup2FAHandler{
		setup2FAUseCase: useCase,
	}
}

// GenerateQRCode godoc
// @Summary Generar código QR para 2FA
// @Description Genera un código QR para configurar Google Authenticator
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param usuario path string true "Nombre de usuario"
// @Success 200 {object} entities.Setup2FAResponse
// @Failure 400 {object} entities.ErrorResponse
// @Failure 500 {object} entities.ErrorResponse
// @Router /users/{usuario}/2fa/setup [post]
func (h *Setup2FAHandler) GenerateQRCode(ctx *gin.Context) {
	usuario := ctx.Param("usuario")
	if usuario == "" {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Usuario requerido",
		})
		return
	}

	qrCode, secret, err := h.setup2FAUseCase.GenerateQRCode(usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, entities.Setup2FAResponse{
		QRCode:  qrCode,
		Secret:  secret,
		Message: "Escanea el código QR con Google Authenticator",
	})
}

// Enable2FA godoc
// @Summary Habilitar 2FA
// @Description Confirma y habilita la autenticación de dos factores
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param usuario path string true "Nombre de usuario"
// @Param request body entities.Setup2FARequest true "Token TOTP"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entities.ErrorResponse
// @Failure 401 {object} entities.ErrorResponse
// @Failure 500 {object} entities.ErrorResponse
// @Router /users/{usuario}/2fa/enable [post]
func (h *Setup2FAHandler) Enable2FA(ctx *gin.Context) {
	usuario := ctx.Param("usuario")
	if usuario == "" {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Usuario requerido",
		})
		return
	}

	var req entities.Setup2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Token TOTP requerido",
		})
		return
	}

	err := h.setup2FAUseCase.Enable2FA(usuario, req.TokenTOTP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "2FA habilitado exitosamente",
	})
}

// Disable2FA godoc
// @Summary Deshabilitar 2FA
// @Description Deshabilita la autenticación de dos factores
// @Tags 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param usuario path string true "Nombre de usuario"
// @Param request body entities.Setup2FARequest true "Token TOTP"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entities.ErrorResponse
// @Failure 401 {object} entities.ErrorResponse
// @Failure 500 {object} entities.ErrorResponse
// @Router /users/{usuario}/2fa/disable [post]
func (h *Setup2FAHandler) Disable2FA(ctx *gin.Context) {
	usuario := ctx.Param("usuario")
	if usuario == "" {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Usuario requerido",
		})
		return
	}

	var req entities.Setup2FARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Token TOTP requerido",
		})
		return
	}

	err := h.setup2FAUseCase.Disable2FA(usuario, req.TokenTOTP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "2FA deshabilitado exitosamente",
	})
}