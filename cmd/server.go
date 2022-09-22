package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true // пропускаем все
	// },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte)
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := &Server{
		clients:       make(map[*websocket.Conn]bool),
		handleMessage: handleMessage,
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(":8080", nil)

	return server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error connect: ", err)
	}
	defer conn.Close()

	log.Println("Client connected:", conn.RemoteAddr().String())

	// сохраняем соединение
	clients[conn] = true
	defer delete(clients, conn)

	for {
		mtype, message, err := conn.ReadMessage() // читаем сообщение
		if err != nil || mtype == websocket.CloseMessage {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error read message: %v", err)
			}
			break
		}

		//go server.MyWriteMessage(message)
		go server.handleMessage(message)
	}
}

// func (server *Server) MyWriteMessage(message []byte, authorMessage *websocket.Conn) {
func (server *Server) MyWriteMessage(message []byte) {
	for conn := range clients {
		// проверка, чтобы не отправлять это сообщение его автору
		// if conn == authorMessage {
		// 	continue
		// }
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}

func main() {

	server := StartServer(messageHandler)

	for {
		server.MyWriteMessage([]byte("Hello"))
	}

}
