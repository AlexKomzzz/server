package handler

import (
	"context"
	"log"
	"net/http"

	chat "github.com/AlexKomzzz/server"
	"github.com/AlexKomzzz/server/pkg/service"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service *service.Service
	clients map[string]map[*websocket.Conn]bool
	ctx     context.Context
}

func NewHandler(service *service.Service, clients map[string]map[*websocket.Conn]bool, ctx context.Context) *Handler {
	return &Handler{
		service: service,
		clients: clients,
		ctx:     ctx,
	}
}

// рассылка сообщений подключенным пользователям
// keyClients - передать chat{idChat} или group{idGroup}
func (h *Handler) sendMessage(msg *chat.Message, keyClients string) {

	// отправим сообщение каждому подключенному клиенту
	for client := range h.clients[keyClients] {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(h.clients[keyClients], client)
		}
	}
}

func (h *Handler) InitRouter() *http.ServeMux {

	// создаем в мапе clients массив для записи подключенных клиентов ОБЩЕГО чата, где ключ будет "all"
	h.clients["all"] = make(map[*websocket.Conn]bool, 0)

	router := http.NewServeMux()

	router.HandleFunc("/test", h.test)

	// Аутентификация и авторизация
	router.HandleFunc("/auth/sign-up", h.signUp)
	router.HandleFunc("/auth/sign-in", h.signIn)

	// Запуск общего чата после авторизации
	// создание общего чата
	//router.Handle("/", h.userIdentity(http.FileServer(http.Dir("./web"))))
	// router.Handle("/chat/start", h.userIdentity(http.StripPrefix("/chat/start", http.FileServer(http.Dir("./web")))))
	router.Handle("/start/", h.userIdentity(http.StripPrefix("/start/", http.FileServer(http.Dir("./web/start/")))))
	router.HandleFunc("/ws", h.WebsocketHandler)

	// создание приватного чата для двоих
	// router.Handle("/chat_two/", h.userIdentity(http.StripPrefix("/chat_two/", http.FileServer(http.Dir("./web/chat_two/")))))
	router.Handle("/chat_two/", h.parseEmailAndIdentity(http.StripPrefix("/chat_two/", http.FileServer(http.Dir("./web/chat_two/")))))
	router.HandleFunc("/chat", h.ChatTwoUser)
	// router.HandleFunc("/chat", h.getChat(h.ChatTwoUser))
	// пример URL http://localhost:8080/chat_two/?email={email_user}

	// Создание чата по email
	// в url должен быть след. фрагмент: ?email=bobik
	// router.HandleFunc("/chat", h.userIdentityHF(h.getChat))

	// создание группового чата
	// в url должен быть след. фрагмент: ?title={title_group}
	router.HandleFunc("/group_chat", h.userIdentityHF(h.getGroup))

	// подключение к групповому чату
	// пример URL http://localhost:8080/chat_group/?idGroup={id_group}
	router.Handle("/chat_group/", h.parseIdGroupAndIdentity(http.StripPrefix("/chat_group/", http.FileServer(http.Dir("./web/chat_group/")))))
	router.HandleFunc("/chat_group", h.ConnGroupChat)

	return router
}
