package infraestructure

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/controllers"
	"apiusuarios/src/usuarios/infraestructure/repositories/mysql"
	"apiusuarios/src/usuarios/infraestructure/routes"

	"github.com/gin-gonic/gin"
)

func InitUser(r *gin.Engine) {
	userRepository := mysql.NewMySql()
	createUserUseCase := application.NewCreateUserUseCase(userRepository)
	loginUserUseCase := application.NewLoginUserUseCase(userRepository)
	tfaUseCase := application.NewTFAUseCase(userRepository)
	changePasswordUseCase := application.NewChangePasswordUseCase(userRepository)
	
	userController := controllers.NewUserController(
		createUserUseCase, 
		loginUserUseCase, 
		tfaUseCase, 
		changePasswordUseCase,
	)
	
	routes.UserRoutes(r, userController)
}