package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	host = "localhost:8080"
)

type Server struct {
	clients map[*websocket.Conn]bool
	//handleMessage func(message []byte)
}

func (server *Server) StartServer() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/ws", server.WebsocketHandler)
	log.Println("сервер запущен на хосту: ", host)
	http.ListenAndServe(host, nil)
}

// func (server *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, _ := template.ParseFiles("templates/index.html")
// 	if err := tmpl.Execute(w, nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func messageHandler(message string) {
	fmt.Println(string(message))
}

func main() {
	server := &Server{
		clients: make(map[*websocket.Conn]bool),
		//handleMessage: handleMessage,
	}
	server.StartServer()

	// for {
	//server.MyWriteMessage([]byte("Hello"))
	// }

}
