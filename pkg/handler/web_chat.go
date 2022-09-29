package handler

import "net/http"

func (h *Handler) StartChat() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/ws", h.webClient.WebsocketHandler)
}
