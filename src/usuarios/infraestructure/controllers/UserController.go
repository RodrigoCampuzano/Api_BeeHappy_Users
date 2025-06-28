package controllers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/handlers"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	createUserUseCase *handlers.CreateUserHandler
}

func NewUserController(createUserUseCase *application.CreateUserUseCase) *UserController {
	createHandler := handlers.NewCreateUserHandler(*createUserUseCase)

	return &UserController{
		createUserUseCase: createHandler,	
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	uc.createUserUseCase.SaveUser(ctx)
}