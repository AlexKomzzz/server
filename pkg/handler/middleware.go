package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// парсинг хедера, определение JWT, определение id
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization") // выделяем из заголовка поле "Authorization"
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	// запись id пользователя в контекст
	c.Set("userId", userId)
}
