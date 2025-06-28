package main

import (
	"apiusuarios/src/core/middleware"
	"apiusuarios/src/usuarios/infraestructure"
	"github.com/gin-gonic/gin"
)

func main() { 
	r := gin.Default()
	r.Use(middleware.MiddlewareCORS())
	infraestructure.InitUser(r)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}