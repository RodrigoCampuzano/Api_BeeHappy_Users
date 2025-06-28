package handlers

import (
    "apiusuarios/src/usuarios/application"
    "net/http"
    "github.com/gin-gonic/gin"
)

type LoginUserHandler struct {
    loginUserUseCase *application.LoginUserUseCase
}

func NewLoginUserHandler(loginUseCase *application.LoginUserUseCase) *LoginUserHandler {
    return &LoginUserHandler{
        loginUserUseCase: loginUseCase,
    }
}

type LoginRequest struct {
    Usuario    string `json:"usuario"`
    Contrasena string `json:"contrasena"`
}

func (h *LoginUserHandler) Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    token, err := h.loginUserUseCase.Execute(req.Usuario, req.Contrasena)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}