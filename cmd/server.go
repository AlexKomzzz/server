package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlexKomzzz/server/pkg/repository"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Server struct {
	db      *sqlx.DB
	clients map[*websocket.Conn]bool

	//handleMessage func(message []byte)
}

func initConfig() error { //Инициализация конфигураций
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func (server *Server) StartServer() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/", server.IndexHandler)
	http.HandleFunc("/ws", server.WebsocketHandler)
	log.Println("сервер запущен на хосту:\t", fmt.Sprint(viper.GetString("host")+viper.GetString("port")))
	http.ListenAndServe(fmt.Sprint(viper.GetString("host")+viper.GetString("port")), nil)
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

	// Инициализируем конфигурации
	if err := initConfig(); err != nil {
		log.Fatalln("error initializing configs: ", err)
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
	})
	if err != nil {
		log.Fatalln("failed to initialize db: ", err)
		return
	}
	defer db.Close()

	server := &Server{
		clients: make(map[*websocket.Conn]bool),
		//handleMessage: handleMessage,
	}
	server.StartServer()

	// for {
	//server.MyWriteMessage([]byte("Hello"))
	// }

}
