package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	chat "github.com/AlexKomzzz/server"
	"github.com/gorilla/websocket"
)

type WebClient struct {
	clients map[*websocket.Conn]bool
	ctx     context.Context
}

func NewWebClient(clients map[*websocket.Conn]bool, ctx context.Context) *WebClient {
	return &WebClient{
		clients: clients,
		ctx:     ctx,
	}
}

// объект сообщения
// type Message struct {
// 	Username string `json:"username"`
// 	Body     string `json:"message"`
// 	Date     string `json:"date"`
// }

var upgrader = websocket.Upgrader{
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true // пропускаем все
	// },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (h *Handler) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	////////
	//upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	////////

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())

	//clnt.ctx = context.WithValue(clnt.ctx, keyName, "Alex")
	//h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, 1)

	// вытащим username пользователя из контекста
	//username := h.webClient.ctx.Value(keyName).(string)

	//userId := h.webClient.ctx.Value(keyId).(int)
	//username, err := h.service.GetUsername(userId)
	username, err := "Alex", nil
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}

	// сохраняем соединение
	h.webClient.clients[conn] = true
	defer delete(h.webClient.clients, conn)

	for {
		var msg *chat.Message
		// читаем сообщение, парсим json в структуру сообщения
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(h.webClient.clients, conn)
			break
		}

		// в сообщение добавим время и username
		msg.Date = time.Now().Format(time.Stamp)
		msg.Username = username

		// сохраняем сообщение в БД
		//h.service.WriteInChat(msg, )

		go h.MyWriteMessage(msg) // отправляем сообщение
		// go messageHandler(message) // выводим сообщение

		//////////////////////////////messageHandler(msg.Message)

	}
}

func (h *Handler) MyWriteMessage(msg *chat.Message) {

	// Grab the next message from the broadcast channel
	// msg := <-broadcast
	// отправим сообщение каждому подключенному клиенту
	for client := range h.webClient.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(h.webClient.clients, client)
		}
	}

	// for conn := range clients {
	// 	conn.WriteMessage(websocket.TextMessage, message)
	// }
}

/*
func (clnt *WebClient) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())
	//var socketClient *ConnectUser = newConnectUser(conn, conn.RemoteAddr().String())
	// сохраняем соединение
	clnt.clients[conn] = true
	defer delete(clnt.clients, conn)
	for {
		var msg Message
		// читаем сообщение, парсим json в структуру сообщения
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clnt.clients, conn)
			break
		}
		// mtype, message, err := conn.ReadMessage() // читаем сообщение
		// if err != nil || mtype == websocket.CloseMessage {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Fatalf("error read message: %v", err)
		// 	}
		// 	break
		// }
		msg.Time = time.Now().Format(time.Stamp)
		go clnt.MyWriteMessage(msg) // отправляем сообщение
		// go messageHandler(message) // выводим сообщение

	}
}
func (clnt *WebClient) MyWriteMessage(msg Message) {
	// Grab the next message from the broadcast channel
	// msg := <-broadcast
	// отправим сообщение каждому подключенному клиенту
	for client := range clnt.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clnt.clients, client)
		}
	}
	// for conn := range clients {
	// 	conn.WriteMessage(websocket.TextMessage, message)
	// }
}
*/
