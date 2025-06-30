package handlers

import (
	"apiusuarios/src/usuarios/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUserHandler struct {
	loginUserUseCase *application.LoginUserUseCase
}

func NewLoginUserHandler(loginUseCase *application.LoginUserUseCase) *LoginUserHandler {
	return &LoginUserHandler{
		loginUserUseCase: loginUseCase,
	}
}

type LoginRequest struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type VerifyLoginRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *LoginUserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := h.loginUserUseCase.Execute(req.Usuario, req.Contrasena)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if result.RequireTwoFactor {
		ctx.JSON(http.StatusOK, gin.H{
			"require_two_factor": true,
			"email":              result.Email,
			"message":            result.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":              result.Token,
		"require_two_factor": false,
		"message":            result.Message,
	})
}

func (h *LoginUserHandler) VerifyLogin(ctx *gin.Context) {
	var req VerifyLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := h.loginUserUseCase.VerifyTwoFactorAndLogin(req.Email, req.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":              result.Token,
		"require_two_factor": result.RequireTwoFactor,
		"message":            result.Message,
	})
}
