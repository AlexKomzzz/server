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

	// Создание приватного чата по idUser2
	// в url должен быть след. фрагмент: ?idUser2={id_user2}
	router.HandleFunc("/new_priv", h.identityAndParseURLHF(h.getPrivChat))

	// подключение к приватному чата с пользователем по его id
	// id пользователя передаем в URL
	// в url должен быть след. фрагмент: ?idUser2={id_user2}
	// пример URL http://localhost:8080/chat_priv?idUser2=3/
	router.Handle("/chat_priv/", h.identityAndParseURL(h.compareIdUser(http.StripPrefix("/chat_priv/", http.FileServer(http.Dir("./web/chat_priv/"))))))
	router.HandleFunc("/chat", h.connPrivChat)
	// router.HandleFunc("/chat", h.getChat(h.ChatTwoUser))

	// создание группового чата
	// в url должен быть след. фрагмент: ?title={title_group}
	router.HandleFunc("/new_group", h.identityAndParseURLHF(h.getGroup))

	// подключение к групповому чату
	// пример URL http://localhost:8080/chat_group/?idGroup={id_group}
	router.Handle("/chat_group/", h.identityAndParseURL(http.StripPrefix("/chat_group/", http.FileServer(http.Dir("./web/chat_group/")))))
	router.HandleFunc("/chat_group", h.ConnGroupChat)

	return router
}
