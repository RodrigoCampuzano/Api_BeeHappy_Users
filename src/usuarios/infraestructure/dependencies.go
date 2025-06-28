package infraestructure

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/controllers"
	"apiusuarios/src/usuarios/infraestructure/repositories/mysql"
	"apiusuarios/src/usuarios/infraestructure/routes"

	"github.com/gin-gonic/gin"
)

func InitUser(r *gin.Engine) {
	ps := mysql.NewMySql()
	createUserHandler := application.NewCreateUserUseCase(ps)
	loginUserHandler := application.NewLoginUserUseCase(ps)
	UserController := controllers.NewUserController(createUserHandler, loginUserHandler)
	routes.UserRoutes(r, UserController)
}