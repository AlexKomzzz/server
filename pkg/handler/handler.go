package handler

import (
	"net/http"

	"github.com/AlexKomzzz/server/pkg/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service   *service.Service
	webClient *WebClient
}

func NewHandler(service *service.Service, webClient *WebClient) *Handler {
	return &Handler{
		service:   service,
		webClient: webClient,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()

	auth := router.PathPrefix("/auth").Methods("POST").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp)
		auth.HandleFunc("/sign-in", h.signIn)
	}

	// открытие websocket
	// вложение Handler в другой Handler для проверки аутентификации
	//router.Handle("/chat/", h.userIdentity(http.StripPrefix("/chat/", http.FileServer(http.Dir("./web")))))
	//router.Handle("/chat/", http.StripPrefix("/chat/", http.FileServer(http.Dir("./web"))))
	router.Handle("/", http.FileServer(http.Dir("./oldWeb")))
	router.HandleFunc("/ws", h.webClient.WebsocketHandler)

	// {
	// 	// chat.GET("/start", h.StartChat)
	// 	chat.Static("/", "./web")
	// 	chat.GET("/ws", h.StartChat)
	// }

	return router
}
