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

// CreateUser godoc
// @Summary Crear nuevo usuario
// @Description Crear un nuevo usuario en el sistema
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param user body entities.User true "Datos del usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entities.ErrorResponse
// @Failure 500 {object} entities.ErrorResponse
// @Router /users [post]
func (uc *CreateUserHandler) SaveUser(ctx *gin.Context) {
	var data *entities.User
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := uc.userCreateUseCase.Execute(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario creado exitosamente",
	})
}