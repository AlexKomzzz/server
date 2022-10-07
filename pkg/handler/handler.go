package handler

import (
	"net/http"

	"github.com/AlexKomzzz/server/pkg/service"
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

func (h *Handler) InitRouter() *http.ServeMux {
	// router := mux.NewRouter()

	// auth := router.PathPrefix("/auth").Methods("POST").Subrouter()
	// {
	// 	auth.HandleFunc("/sign-up", h.signUp)
	// 	auth.HandleFunc("/sign-in", h.signIn)
	// }

	// открытие websocket
	// вложение Handler в другой Handler для проверки аутентификации
	//router.Handle("/chat/", h.userIdentity(http.StripPrefix("/chat/", http.FileServer(http.Dir("./web")))))
	//router.Handle("/chat/", http.StripPrefix("/chat/", http.FileServer(http.Dir("./web"))))
	//router.Handle("/", h.userIdentity(http.FileServer(http.Dir("./web"))))

	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	// // http.Handle("/", http.FileServer(http.Dir("./web")))
	// // router.Handle("/", http.FileServer(http.Dir("./web")))
	// router.HandleFunc("/ws", h.webClient.WebsocketHandler)

	// {
	// 	// chat.GET("/start", h.StartChat)
	// 	chat.Static("/", "./web")
	// 	chat.GET("/ws", h.StartChat)
	// }

	router := http.NewServeMux()

	router.HandleFunc("/test", h.test)

	// Аутентификация и авторизация
	router.HandleFunc("/auth/sign-up", h.signUp)
	router.HandleFunc("/auth/sign-in", h.signIn)

	// Запуск чата после авторизации
	router.Handle("/", h.userIdentity(http.FileServer(http.Dir("./web"))))
	router.HandleFunc("/ws", h.webClient.WebsocketHandler)

	// Создание чата по email
	router.HandleFunc("/chat/:email", h.getChat)

	return router
}
