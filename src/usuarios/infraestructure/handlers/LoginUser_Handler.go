package handlers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/domain/entities"
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

// Login godoc
// @Summary Iniciar sesión
// @Description Autenticar usuario con contraseña y opcionalmente token 2FA
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body entities.LoginRequest true "Credenciales de login"
// @Success 200 {object} entities.LoginResponse
// @Failure 400 {object} entities.ErrorResponse
// @Failure 401 {object} entities.ErrorResponse
// @Router /users/login [post]
func (h *LoginUserHandler) Login(ctx *gin.Context) {
	var req entities.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Datos inválidos",
		})
		return
	}

	token, err := h.loginUserUseCase.Execute(req.Usuario, req.Contrasena, req.TokenTOTP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, entities.LoginResponse{
		Token:   token,
		Message: "Login exitoso",
	})
}