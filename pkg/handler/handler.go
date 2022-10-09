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
	// Запуск общего чата после авторизации
	router.Handle("/start/", h.userIdentity(http.StripPrefix("/start/", http.FileServer(http.Dir("./web/start/")))))
	// создание общего чата
	router.HandleFunc("/ws", h.WebsocketHandler)

	// создание приватного чата для двоих
	// router.Handle("/chat_two/", h.userIdentity(http.StripPrefix("/chat_two/", http.FileServer(http.Dir("./web/chat_two/")))))
	router.Handle("/chat_two/", h.parseEmailAndIdentity(http.StripPrefix("/chat_two/", http.FileServer(http.Dir("./web/chat_two/")))))

	router.HandleFunc("/chat", h.getChat(h.ChatTwoUser))
	//router.HandleFunc("/chat", h.WebsocketHandler)

	// Создание чата по email
	// в url должен быть след. фрагмент: ?email=bobik
	// router.HandleFunc("/chat", h.userIdentityHF(h.getChat))

	return router
}
