package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zardan4/todo-app-gin"
)

// реєстрація
func (h *Handler) signUp(c *gin.Context) {
	// проходиться валідація по тегу binding required

	var input todo.User

	if err := c.BindJSON(&input); err != nil { // парсимо отриманий від користувача JSON
		newErrorResponse(c, http.StatusBadRequest, err.Error()) // помилка
		return
	}

	// передаємо дані на шар нижче
	id, err := h.services.Authorization.CreateUser(input) // створюємо користувача
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // помилка
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	}) // відправляємо айді створеного юзера
}

// аутентифікація
type signInInput struct { // структура, яка буде приходити від користувача при аутентифікації
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// передаємо дані на шар нижче
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // генеруємо токен в бізнес-логіці
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) // помилка
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	}) // відправляємо айді створеного юзера

}
