package main

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "apiusuarios/docs"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.MiddlewareCORS())
	infraestructure.InitUser(r)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
