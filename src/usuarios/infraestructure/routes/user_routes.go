package routes

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userController *controllers.UserController) {
	usergroup := r.Group("/users")
	{
		// Rutas existentes
		usergroup.POST("/", userController.CreateUser)
		usergroup.POST("/login", userController.Login) // Login sin 2FA (mantenido para compatibilidad)
		
		// Nuevas rutas para 2FA
		usergroup.POST("/login/2fa", userController.LoginWithTFA)           // Paso 1: Validar credenciales y enviar código
		usergroup.POST("/login/verify", userController.VerifyLoginTFA)      // Paso 2: Verificar código y obtener token
		
		// Rutas para gestión de códigos 2FA
		usergroup.POST("/2fa/send-login-code", userController.SendLoginCode)
		usergroup.POST("/2fa/verify-login-code", userController.VerifyLoginCode)
		
		// Rutas para cambio de contraseña con 2FA
		usergroup.POST("/password/send-code", userController.SendPasswordChangeCode)
		usergroup.POST("/password/change", userController.ChangePassword)

		// Rutas protegidas
		protected := usergroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Ejemplo: protected.GET("/perfil", userService.Perfil)
			// Agrega aquí las rutas que quieras proteger
		}
	}
}
