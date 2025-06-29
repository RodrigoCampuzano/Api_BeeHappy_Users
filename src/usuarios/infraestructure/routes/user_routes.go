package routes

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userController *controllers.UserController) {
	userGroup := r.Group("/users")
	{
		// Rutas de usuarios
		userGroup.POST("/", userController.CreateUser)
		userGroup.POST("/login", userController.LoginUser)
		
		// Rutas de 2FA
		userGroup.POST("/:usuario/2fa/setup", userController.GenerateQRCode)
		userGroup.POST("/:usuario/2fa/enable", userController.Enable2FA)
		userGroup.POST("/:usuario/2fa/disable", userController.Disable2FA)

		// Rutas protegidas
		protected := usergroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Ejemplo: protected.GET("/perfil", userService.Perfil)
			// Agrega aquí las rutas que quieras proteger
		}
	}
}