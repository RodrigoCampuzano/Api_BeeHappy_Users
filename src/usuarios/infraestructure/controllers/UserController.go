package controllers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/infraestructure/handlers"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	createUserUseCase *handlers.CreateUserHandler
	loginUserUseCase *handlers.LoginUserHandler
}

func NewUserController(createUserUseCase *application.CreateUserUseCase, loginUserUseCase *application.LoginUserUseCase) *UserController {
	createHandler := handlers.NewCreateUserHandler(*createUserUseCase)
	loginHandler := handlers.NewLoginUserHandler(*&loginUserUseCase)

	return &UserController{
		createUserUseCase: createHandler,	
		loginUserUseCase: loginHandler,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	uc.createUserUseCase.SaveUser(ctx)
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	uc.loginUserUseCase.Login(ctx)
}