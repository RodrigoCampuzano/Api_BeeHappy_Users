// Package main API de Usuarios
package main

import (
	"apiusuarios/docs"
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"


	_ "apiusuarios/docs" // Importa los docs generados
)

// @title           API de Usuarios
// @version         1.0
// @description     API para gestión de usuarios con autenticación en dos pasos y recuperación de contraseña.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Introduce el token con el formato: Bearer <token>

func main() {

	r := gin.Default()

	// Middleware CORS
	r.Use(middleware.MiddlewareCORS())
	// Configuración de Swagger
	config := &ginSwagger.Config{
		URL:                      "http://localhost:8080/swagger/doc.json", // La URL donde se sirve el JSON de Swagger
		DeepLinking:              true,
		DocExpansion:             "list",
		DefaultModelsExpandDepth: 1,
	}
	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(config, swaggerFiles.Handler))

	// Grupo base para la API v1
	api := r.Group("/api/v1")
	{
		// Inicializar rutas
		infraestructure.InitUser(api)
	}

	log.Fatal(r.Run(":8080"))
}
