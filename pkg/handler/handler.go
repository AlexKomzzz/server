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
	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
	router.Handle("/", http.FileServer(http.Dir("./web")))
	router.HandleFunc("/ws", h.webClient.WebsocketHandler)
	// mux.Static("/", "./web")
	// mux.GET("/ws", h.webClient.WebsocketHandler)
	//chat := mux.Group("/chat", h.userIdentity)
	// chat := mux.Group("/chat")

	// {
	// 	// chat.GET("/start", h.StartChat)
	// 	chat.Static("/", "./web")
	// 	chat.GET("/ws", h.StartChat)
	// }

	return router
}
