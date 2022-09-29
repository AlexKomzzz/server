package handler

import (
	"fmt"
	"net/http"

	chat "github.com/AlexKomzzz/server"
	"github.com/gin-gonic/gin"
)

// Обработчик для регистрации пользователя
func (h *Handler) signUp(c *gin.Context) {
	var user chat.User

	// парсим тело запроса в структуру пользователя
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	// по данным пользователя заносим в БД и получаем id
	id, err := h.service.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})

}

type InUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Обработчик для аутентификации пользователя
func (h *Handler) signIn(c *gin.Context) {
	var user InUser

	// парсим тело запроса в структуру пользователя
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("invalid input body: %v", err))
		return
	}

	// по данным пользователя заносим в БД и получаем id
	token, err := h.service.GenerateToken(user.Email, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
