package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	chat "github.com/AlexKomzzz/server"
	"github.com/gorilla/websocket"
)

type HistoryResp struct {
	History []chat.Message
}

// создание чата с другим пользователем по его email
/*func (h *Handler) getChat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "invalid method: no GET", http.StatusBadRequest)
			return
		}

		//// для проверки без идентификации
		// h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, 1)
		////

		// var historyChat []chat.Message

		// // вытащим id пользователя из контекста
		// idUser := h.webClient.ctx.Value(keyId).(int)

		// получение истории чата с пользователем
		// historyChat, err = h.service.GetChat(idUser, emailUser2)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// // http ответ
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(&HistoryResp{
		// 	History: historyChat,
		// })

		next(w, r)
	}
}*/

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (h *Handler) ChatTwoUser(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())

	////////////////////////
	// username, err := "Alex", nil
	// if err != nil {
	// 	log.Fatalln("error: не получен username по id: ", err)
	// }
	////////////////////////

	// сохраняем соединение
	clients := make(map[*websocket.Conn]bool)
	clients[conn] = true
	defer delete(clients, conn)

	// получение id текущего пользователя из контекста
	idUser1 := h.webClient.ctx.Value(keyId).(int)
	log.Println("id = ", idUser1)
	// получение email пользователя, с которым создаем чат, из контекста
	emailUser2 := h.webClient.ctx.Value(keyEmail).(string)

	// получение username по id
	username, err := h.service.GetUsername(idUser1)
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}

	// получение истории чата из БД
	historyChat, err := h.service.GetChat(idUser1, emailUser2)
	if err != nil {
		log.Fatalln("error: не получена история чата: ", err)
	}

	// historyChat := []*chat.Message{{
	// 	Date:     "2004-10-19 10:23:54",
	// 	Username: "Alex",
	// 	Body:     "Hello",
	// }}

	// передача истории клиентам
	if len(historyChat) > 0 {
		for _, msg := range historyChat {
			msg.Date = strings.Replace(strings.Replace(msg.Date, "T", " ", 1), "Z", "       ", 1)
			h.sendMessage(msg, clients) // возможна блокировка
		}
	}

	for {
		var msg *chat.Message
		// читаем сообщение, парсим json в структуру сообщения
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}

		// в сообщение добавим время и username
		msg.Date = time.Now().Format("2006-01-02 15:04:05")
		msg.Username = username

		// сохраняем сообщение в БД
		err = h.service.WriteInChat(msg, idUser1, emailUser2)
		if err != nil {
			log.Fatalln("error: сообщение не сохранено: ", err)
		}

		// отправляем сообщение
		h.sendMessage(msg, clients)
	}
}

func (h *Handler) sendMessage(msg *chat.Message, clients map[*websocket.Conn]bool) {

	// отправим сообщение каждому подключенному клиенту
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
