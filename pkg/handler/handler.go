package handler

import (
	"github.com/AlexKomzzz/server/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRouter() {
	mux := gin.New()
	mux.Group("/auth")
	{
		mux.POST("/sign-up", h.signUp)
		mux.POST("/sign-in", h.signIn)
	}
}
