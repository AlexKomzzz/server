package handler

import (
	"github.com/AlexKomzzz/server/pkg/service"
	"github.com/AlexKomzzz/server/pkg/webclient"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *service.Service
	webClient *webclient.WebClient
}

func NewHandler(service *service.Service, webClient *webclient.WebClient) *Handler {
	return &Handler{
		service:   service,
		webClient: webClient,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	mux := gin.New()
	auth := mux.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	// открытие websocket
	mux.Static("/", "./web")
	mux.GET("/ws", h.webClient.WebsocketHandler)
	//chat := mux.Group("/chat", h.userIdentity)
	// chat := mux.Group("/chat")

	// {
	// 	// chat.GET("/start", h.StartChat)
	// 	chat.Static("/", "./web")
	// 	chat.GET("/ws", h.StartChat)
	// }

	return mux
}
