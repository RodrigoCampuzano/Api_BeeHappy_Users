package infraestructure

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/controllers"
	"apiusuarios/src/usuarios/infraestructure/repositories/mysql"
	"apiusuarios/src/usuarios/infraestructure/routes"

	"github.com/gin-gonic/gin"
)

func InitUser(r *gin.Engine) {
	// Repositorio
	userRepo := mysql.NewMySql()
	
	// Casos de uso
	createUserUseCase := application.NewCreateUserUseCase(userRepo)
	loginUserUseCase := application.NewLoginUserUseCase(userRepo)
	setup2FAUseCase := application.NewSetup2FAUseCase(userRepo)
	
	// Controlador
	userController := controllers.NewUserController(createUserUseCase, loginUserUseCase, setup2FAUseCase)
	
	// Rutas
	routes.UserRoutes(r, userController)
}