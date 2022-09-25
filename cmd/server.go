package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectUser struct {
	Websocket *websocket.Conn
	ClientIP  string
}

func newConnectUser(ws *websocket.Conn, clientIP string) *ConnectUser {
	return &ConnectUser{
		Websocket: ws,
		ClientIP:  clientIP,
	}
}

var clients = make(map[ConnectUser]int)

//var clients map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true // пропускаем все
	// },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	clients map[*websocket.Conn]bool
	//handleMessage func(message []byte)
}

func (server *Server) StartServer(handleMessage func(message []byte)) {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/ws", server.WebsocketHandler)
	http.ListenAndServe("localhost:8080", nil)
}

// func (server *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, _ := template.ParseFiles("templates/index.html")
// 	if err := tmpl.Execute(w, nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()

	log.Println("Client connected:", conn.RemoteAddr().String())
	var socketClient *ConnectUser = newConnectUser(conn, conn.RemoteAddr().String())
	// сохраняем соединение
	//clients[conn] = true
	clients[*socketClient] = 0
	defer delete(clients, *socketClient)

	for {
		mtype, message, err := conn.ReadMessage() // читаем сообщение
		if err != nil || mtype == websocket.CloseMessage {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Fatalf("error read message: %v", err)
			}
			break
		}

		timeMessage := time.Now().Format(time.Stamp)
		message = []byte(fmt.Sprintf("%s %s\t\t%s", socketClient.ClientIP, timeMessage, string(message)))

		go server.MyWriteMessage(message) // отправляем сообщение
		// go messageHandler(message) // выводим сообщение
		messageHandler(message)
	}
}

func (server *Server) MyWriteMessage(message []byte) {
	for user := range clients {
		user.Websocket.WriteMessage(websocket.TextMessage, message)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}

func main() {
	server := &Server{
		clients: make(map[*websocket.Conn]bool),
		//handleMessage: handleMessage,
	}
	server.StartServer(messageHandler)

	// for {
	//server.MyWriteMessage([]byte("Hello"))
	// }

}
