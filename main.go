package main

import (
	"apiusuarios/docs"
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Usuarios con 2FA
// @version 1.0
// @description API para gestión de usuarios con autenticación de dos factores usando Google Authenticator
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Configurar información de Swagger
	docs.SwaggerInfo.Title = "API de Usuarios con 2FA"
	docs.SwaggerInfo.Description = "API para gestión de usuarios con autenticación de dos factores"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	r.Use(middleware.MiddlewareCORS())
	
	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	infraestructure.InitUser(r)
	
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}