package controllers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/handlers"

	"github.com/gin-gonic/gin"
)

// UserController maneja las rutas relacionadas con usuarios
type UserController struct {
	createUserUseCase    *handlers.CreateUserHandler
	loginUserUseCase     *handlers.LoginUserHandler
	twoFactorAuthHandler *handlers.TwoFactorAuthHandler
	passwordResetHandler *handlers.PasswordResetHandler
}

func NewUserController(
	createUserUseCase *application.CreateUserUseCase,
	loginUserUseCase *application.LoginUserUseCase,
	twoFactorAuthUseCase *application.TwoFactorAuthUseCase,
	passwordResetUseCase *application.PasswordResetUseCase,
) *UserController {
	createHandler := handlers.NewCreateUserHandler(*createUserUseCase)
	loginHandler := handlers.NewLoginUserHandler(loginUserUseCase)
	twoFactorHandler := handlers.NewTwoFactorAuthHandler(twoFactorAuthUseCase)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetUseCase)

	return &UserController{
		createUserUseCase:    createHandler,
		loginUserUseCase:     loginHandler,
		twoFactorAuthHandler: twoFactorHandler,
		passwordResetHandler: passwordResetHandler,
	}
}

// @Summary     Crear nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags        usuarios
// @Accept      json
// @Produce     json
// @Param       user body entities.SwaggerCreateUserRequest true "Datos del usuario"
// @Success     200 {object} entities.SwaggerCreateUserResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Router      /users [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	uc.createUserUseCase.SaveUser(ctx)
}

// @Summary     Iniciar sesión
// @Description Inicia sesión con usuario y contraseña
// @Tags        autenticación
// @Accept      json
// @Produce     json
// @Param       credentials body entities.SwaggerLoginRequest true "Credenciales"
// @Success     200 {object} entities.SwaggerLoginResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Router      /users/login [post]
func (uc *UserController) LoginUser(ctx *gin.Context) {
	uc.loginUserUseCase.Login(ctx)
}

// @Summary     Verificar código de dos pasos
// @Description Verifica el código de autenticación en dos pasos
// @Tags        autenticación
// @Accept      json
// @Produce     json
// @Param       verification body entities.SwaggerVerifyCodeRequest true "Datos de verificación"
// @Success     200 {object} entities.SwaggerVerifyCodeResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Router      /users/login/verify [post]
func (uc *UserController) VerifyLogin(ctx *gin.Context) {
	uc.loginUserUseCase.VerifyLogin(ctx)
}

// @Summary     Solicitar restablecimiento de contraseña
// @Description Envía un código de verificación al correo para restablecer la contraseña
// @Tags        contraseña
// @Accept      json
// @Produce     json
// @Param       request body entities.SwaggerPasswordResetRequest true "Correo electrónico"
// @Success     200 {object} entities.SwaggerPasswordResetResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Router      /users/password/reset/request [post]
func (uc *UserController) RequestPasswordReset(ctx *gin.Context) {
	uc.passwordResetHandler.RequestReset(ctx)
}

// @Summary     Restablecer contraseña
// @Description Restablece la contraseña usando el código de verificación
// @Tags        contraseña
// @Accept      json
// @Produce     json
// @Param       reset body entities.SwaggerResetPasswordRequest true "Datos de restablecimiento"
// @Success     200 {object} entities.SwaggerResetPasswordResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Router      /users/password/reset [post]
func (uc *UserController) ResetPassword(ctx *gin.Context) {
	uc.passwordResetHandler.ResetPassword(ctx)
}

// @Summary     Solicitar cambio de contraseña
// @Description Envía un código de verificación para cambiar la contraseña
// @Tags        contraseña
// @Accept      json
// @Produce     json
// @Param       change body entities.SwaggerChangePasswordRequest true "Datos de cambio"
// @Success     200 {object} entities.SwaggerChangePasswordResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Security    ApiKeyAuth
// @Router      /users/profile/password/change/request [post]
func (uc *UserController) RequestChangePassword(ctx *gin.Context) {
	uc.passwordResetHandler.RequestChangePassword(ctx)
}

// @Summary     Cambiar contraseña
// @Description Cambia la contraseña usando el código de verificación
// @Tags        contraseña
// @Accept      json
// @Produce     json
// @Param       change body entities.SwaggerResetPasswordRequest true "Datos de cambio"
// @Success     200 {object} entities.SwaggerResetPasswordResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Security    ApiKeyAuth
// @Router      /users/profile/password/change [post]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
	uc.passwordResetHandler.ChangePassword(ctx)
}

// @Summary     Activar/desactivar verificación en dos pasos
// @Description Activa o desactiva la verificación en dos pasos para el usuario
// @Tags        seguridad
// @Accept      json
// @Produce     json
// @Param       toggle body entities.SwaggerToggleTwoFactorRequest true "Estado de verificación"
// @Success     200 {object} entities.SwaggerToggleTwoFactorResponse
// @Failure     400 {object} entities.SwaggerErrorResponse
// @Failure     401 {object} entities.SwaggerErrorResponse
// @Security    ApiKeyAuth
// @Router      /users/profile/2fa/toggle [post]
func (uc *UserController) ToggleTwoFactor(ctx *gin.Context) {
	uc.twoFactorAuthHandler.ToggleTwoFactor(ctx)
}
