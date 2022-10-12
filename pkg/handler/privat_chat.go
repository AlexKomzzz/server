package handler

import (
	"fmt"
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

// создание приватного чата по id пользователей
// возвращает id созданного чата
func (h *Handler) getChat(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// вытащим id пользователя из контекста
	idUser1 := h.ctx.Value(keyId).(int)
	idUser2 := h.ctx.Value(keyIdUser2).(int)

	// создание чата с пользователем
	idChat, err := h.service.CreatePrivChat(idUser1, idUser2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// создаем в мапе clients клиентов мапу для записи подключенных клиентов, где ключ будет "chat{idChat}"
	h.clients[fmt.Sprintf("chat%d", idChat)] = make(map[*websocket.Conn]bool, 0)

	// http ответ
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\n\t\"idChat\": \"%d\"\n}", idChat)))
}

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

	// получение id текущего пользователя из контекста
	idUser1 := h.ctx.Value(keyId).(int)
	log.Println("id = ", idUser1)
	// получение idUser2 пользователя, с которым создаем чат, из контекста
	idUser2 := h.ctx.Value(keyIdUser2).(int)

	// получение id приватного чата, к которому осуществилось подключение
	idChat, err := h.service.GetIdPrivChat(idUser1, idUser2)
	if err != nil {
		log.Fatalln("error: не получен idChat: ", err)
	}

	keyClients := fmt.Sprintf("chat%d", idChat)

	// сохраняем соединение
	h.clients[keyClients][conn] = true
	defer delete(h.clients[keyClients], conn)

	// получение username по id
	username, err := h.service.GetUsername(idUser1)
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}

	// получение истории чата из БД
	historyChat, err := h.service.GetPrivChat(idUser1, idUser2)
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
			h.sendMessage(msg, keyClients) // возможна блокировка
		}
	}

	for {
		var msg *chat.Message
		// читаем сообщение, парсим json в структуру сообщения
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(h.clients[keyClients], conn)
			break
		}

		// в сообщение добавим время и username
		msg.Date = time.Now().Format("2006-01-02 15:04:05")
		msg.Username = username

		// сохраняем сообщение в БД
		err = h.service.WriteInPrivChat(msg, idChat)
		if err != nil {
			log.Fatalln("error: сообщение не сохранено: ", err)
		}

		// отправляем сообщение
		go h.sendMessage(msg, keyClients)
	}
}
