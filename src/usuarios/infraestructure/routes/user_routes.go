package routes

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userService *controllers.UserController) {
	usergroup := r.Group("/users")
	{
		usergroup.POST("/", userService.CreateUser)
		usergroup.POST("/login", userService.LoginUser)

		// Rutas protegidas
		protected := usergroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Ejemplo: protected.GET("/perfil", userService.Perfil)
			// Agrega aquí las rutas que quieras proteger
		}
	}
}
