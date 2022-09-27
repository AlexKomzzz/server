package handler

import (
	"net/http"

	chat "github.com/AlexKomzzz/server"
	"github.com/gin-gonic/gin"
)

// Обработчик для регистрации пользователя
func signUp(c *gin.Context) {
	var user *chat.User

	// парсим тело запроса в структуру пользователя
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	// по данным пользователя заносим в БД

}

// Обработчик для аутентификации пользователя
func signIn(c *gin.Context) {}
