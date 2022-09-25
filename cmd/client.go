package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// объект сообщения
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

//var clients map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true // пропускаем все
	// },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (srv *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()

	log.Println("Client connected:", conn.RemoteAddr().String())
	//var socketClient *ConnectUser = newConnectUser(conn, conn.RemoteAddr().String())
	// сохраняем соединение
	srv.clients[conn] = true
	defer delete(srv.clients, conn)

	for {
		var msg Message
		// читаем сообщение, парсим json в структуру сообщения
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(srv.clients, conn)
			break
		}
		// Send the newly received message to the broadcast channel
		//broadcast <- msg

		// mtype, message, err := conn.ReadMessage() // читаем сообщение
		// if err != nil || mtype == websocket.CloseMessage {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Fatalf("error read message: %v", err)
		// 	}
		// 	break
		// }

		msg.Time = time.Now().Format(time.Stamp)

		go srv.MyWriteMessage(msg) // отправляем сообщение
		// go messageHandler(message) // выводим сообщение
		messageHandler(msg.Message)
	}
}

func (srv *Server) MyWriteMessage(msg Message) {

	// Grab the next message from the broadcast channel
	// msg := <-broadcast
	// отправим сообщение каждому подключенному клиенту
	for client := range srv.clients {
		log.Println("write")
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(srv.clients, client)
		}
	}

	// for conn := range clients {
	// 	conn.WriteMessage(websocket.TextMessage, message)
	// }
}
