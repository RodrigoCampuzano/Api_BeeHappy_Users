package routes

import (
	"apiusuarios/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userService *controllers.UserController) {
	usergroup := r.Group("/users")
	{
		usergroup.POST("/", userService.CreateUser)
	}
}
