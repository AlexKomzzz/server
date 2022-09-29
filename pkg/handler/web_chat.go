package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StartChat(c *gin.Context) {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/ws", h.webClient.WebsocketHandler)
}
