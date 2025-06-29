package routes

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r gin.IRouter, userService *controllers.UserController) {
	// Grupo de usuarios
	usergroup := r.(*gin.RouterGroup).Group("/users")
	{
		// Rutas públicas
		usergroup.POST("/", userService.CreateUser)
		usergroup.POST("/login", userService.LoginUser)
		usergroup.POST("/login/verify", userService.VerifyLogin)

		// Rutas de recuperación de contraseña
		usergroup.POST("/password/reset/request", userService.RequestPasswordReset)
		usergroup.POST("/password/reset", userService.ResetPassword)

		// Rutas protegidas
		protected := usergroup.Group("/profile")
		protected.Use(middleware.AuthMiddleware())
		{
			// Perfil y seguridad
			protected.POST("/password/change/request", userService.RequestChangePassword)
			protected.POST("/password/change", userService.ChangePassword)
			protected.POST("/2fa/toggle", userService.ToggleTwoFactor)
		}
	}
}
