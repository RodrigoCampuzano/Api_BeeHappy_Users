package infraestructure

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/controllers"
	"apiusuarios/src/usuarios/infraestructure/repositories/mysql"
	"apiusuarios/src/usuarios/infraestructure/routes"

	"github.com/gin-gonic/gin"
)

func InitUser(r gin.IRouter) {
	ps := mysql.NewMySql()
	createUserHandler := application.NewCreateUserUseCase(ps)
	loginUserHandler, err := application.NewLoginUserUseCase(ps)
	if err != nil {
		panic(err)
	}

	twoFactorAuthHandler, err := application.NewTwoFactorAuthUseCase(ps)
	if err != nil {
		panic(err)
	}

	passwordResetHandler, err := application.NewPasswordResetUseCase(ps)
	if err != nil {
		panic(err)
	}

	UserController := controllers.NewUserController(createUserHandler, loginUserHandler, twoFactorAuthHandler, passwordResetHandler)
	routes.UserRoutes(r, UserController)
}
