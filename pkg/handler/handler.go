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

	router := http.NewServeMux()

	router.HandleFunc("/test", h.test)

	// Аутентификация и авторизация
	router.HandleFunc("/auth/sign-up", h.signUp)
	router.HandleFunc("/auth/sign-in", h.signIn)

	//router.Handle("/", h.userIdentity(http.FileServer(http.Dir("./web"))))
	// router.Handle("/chat/start", h.userIdentity(http.StripPrefix("/chat/start", http.FileServer(http.Dir("./web")))))
	// Запуск чата после авторизации. добавить start в URL
	router.Handle("/", h.userIdentity(http.StripPrefix("/", http.FileServer(http.Dir("./web/")))))
	router.HandleFunc("/ws", h.WebsocketHandler)

	// Создание чата по email
	// в url должен быть след. фрагмент: ?email={emailUser}
	router.HandleFunc("/chat", h.userIdentityHF(h.getChat))

	return router
}
