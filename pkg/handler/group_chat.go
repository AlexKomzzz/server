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

// создание группового чата
// метод POST
func (h *Handler) getGroup(w http.ResponseWriter, r *http.Request) {

	// проверка метода
	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// получение id текущего пользователя и название группы из контекста

	idUser := h.ctx.Value(keyId).(int)
	log.Printf("idUser = %d", idUser)
	title_group := h.ctx.Value(keyTitle).(string)

	// создаем групповой чат в БД и получаем его id
	idGroup, err := h.service.CreateGroup(title_group, idUser)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("idGroup = %d", idGroup)

	// создаем в мапе clients клиентов массив для записи подключенных клиентов, где ключ будет "group{idGroup}"
	h.clients[fmt.Sprintf("group%d", idGroup)] = make(map[*websocket.Conn]bool)

	// отправим клиенту id группового чата
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\n\t\"idGroupChat\": \"%d\"\n}", idGroup)))
}

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (h *Handler) ConnGroupChat(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())

	// получение id текущего пользователя из контекста
	idUser := h.ctx.Value(keyId).(int)

	// получение id группового чата, к которому осуществилось подключение
	idGroup := h.ctx.Value(keyIdGroup).(int)

	keyClients := fmt.Sprintf("group%d", idGroup)

	// // сохраняем соединение в мапу слиентов
	h.clients[keyClients][conn] = true
	defer delete(h.clients[keyClients], conn)

	// получение username по id
	username, err := h.service.GetUsername(idUser)
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}

	log.Printf("username = %s", username)

	// получение истории группового чата из БД
	historyChat, err := h.service.GetGroup(idGroup)
	if err != nil {
		log.Fatalln("error: не получена история чата: ", err)
	}

	// передача истории клиенту
	if len(historyChat) > 0 {
		for _, msg := range historyChat {
			msg.Date = strings.Replace(strings.Replace(msg.Date, "T", " ", 1), "Z", "       ", 1)

			// отправка сообщения из истории группового чата клиенту
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: при отправке истории группового чата клиенту %v", err)
				conn.Close()
				delete(h.clients[keyClients], conn)
			}
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
		err = h.service.WriteInGroup(msg, idGroup)
		if err != nil {
			log.Fatalln("error: сообщение не сохранено: ", err)
		}

		// отправляем сообщение
		go h.sendMessage(msg, keyClients)
	}
}
