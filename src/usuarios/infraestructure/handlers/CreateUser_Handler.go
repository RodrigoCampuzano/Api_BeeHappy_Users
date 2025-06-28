package handlers

import (
	"apiusuarios/src/usuarios/application"
	"apiusuarios/src/usuarios/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserHandler struct {
	userCreateUseCase *application.CreateUserUseCase
}

func NewCreateUserHandler(userService application.CreateUserUseCase) *CreateUserHandler {
	return &CreateUserHandler{
		userCreateUseCase: &userService,
	}
}

func (uc *CreateUserHandler) SaveUser(ctx *gin.Context) {
	var data *entities.User
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.userCreateUseCase.Execute(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data saved"})
}