// src/usuarios/infraestructure/controllers/UserController.go (ACTUALIZADO)
package controllers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/handlers"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	createUserHandler    *handlers.CreateUserHandler
	loginUserHandler     *handlers.LoginUserHandler
	tfaHandler          *handlers.TFAHandler
}

func NewUserController(
	createUserUseCase *application.CreateUserUseCase, 
	loginUserUseCase *application.LoginUserUseCase,
	tfaUseCase *application.TFAUseCase,
	changePasswordUseCase *application.ChangePasswordUseCase,
) *UserController {
	createHandler := handlers.NewCreateUserHandler(*createUserUseCase)
	loginHandler := handlers.NewLoginUserHandler(loginUserUseCase, tfaUseCase)
	tfaHandler := handlers.NewTFAHandler(tfaUseCase, changePasswordUseCase)

	return &UserController{
		createUserHandler: createHandler,	
		loginUserHandler:  loginHandler,
		tfaHandler:       tfaHandler,
	}
}

// Métodos existentes
// CreateUser godoc
// @Summary Crear usuario
// @Description Crea un nuevo usuario
// @Tags usuarios
// @Accept json
// @Produce json
// @Param user body object true "Datos del usuario"
// @Success 201 {object} object
// @Failure 400 {object} object
// @Router /usuarios [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	uc.createUserHandler.SaveUser(ctx)
}

// Login godoc
// @Summary Login usuario
// @Description Inicia sesión de usuario
// @Tags usuarios
// @Accept json
// @Produce json
// @Param credentials body object true "Credenciales"
// @Success 200 {object} object
// @Failure 401 {object} object
// @Router /usuarios/login [post]
func (uc *UserController) Login(ctx *gin.Context) {
	uc.loginUserHandler.Login(ctx)
}

// Nuevos métodos para 2FA
func (uc *UserController) LoginWithTFA(ctx *gin.Context) {
	uc.loginUserHandler.LoginWithTFA(ctx)
}

func (uc *UserController) VerifyLoginTFA(ctx *gin.Context) {
	uc.loginUserHandler.VerifyLoginTFA(ctx)
}

func (uc *UserController) SendLoginCode(ctx *gin.Context) {
	uc.tfaHandler.SendLoginCode(ctx)
}

func (uc *UserController) VerifyLoginCode(ctx *gin.Context) {
	uc.tfaHandler.VerifyLoginCode(ctx)
}

func (uc *UserController) SendPasswordChangeCode(ctx *gin.Context) {
	uc.tfaHandler.SendPasswordChangeCode(ctx)
}

func (uc *UserController) ChangePassword(ctx *gin.Context) {
	uc.tfaHandler.ChangePassword(ctx)
}