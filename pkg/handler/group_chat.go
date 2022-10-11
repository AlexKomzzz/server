package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	// получение id текущего пользователя из контекста
	idUser := h.ctx.Value(keyId).(int)

	// выделим title_group из url
	// получим мапу из параметров указанных в url с помощью "?"
	set, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
		return
	}
	var title_group string
	// // из мапы вытащим значение email
	if _, ok := set["title"]; ok {
		title_group = set["title"][0]
	}

	// если title_group пустое, значит оно не передано в URL
	if title_group == "" {
		log.Println("Значение title_group не передано в URL")
		http.Error(w, "Значение title_group не передано в URL", http.StatusBadRequest)
		return
	}

	// создаем групповой чат в БД и получаем его id
	idGroup, err := h.service.CreateGroup(title_group, idUser)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// создаем в мапе clients клиентов массив для записи подключенных клиентов, где ключ будет "group{idGroup}"
	h.clients[fmt.Sprintf("group%d", idGroup)] = make(map[*websocket.Conn]bool, 0)

	// отправим клиенту id группового чата
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\n\t\"idGroup\": \"%d\"\n}", idGroup)))
}

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (h *Handler) ConnGroupChat(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())

	// получение id группового чата, к которому осуществилось подключение
	idGroup := h.ctx.Value(keyIdGroup).(int)

	keyClients := fmt.Sprintf("group%d", idGroup)

	// // сохраняем соединение в мапу слиентов
	h.clients[keyClients][conn] = true
	defer delete(h.clients[keyClients], conn)

	// получение id текущего пользователя из контекста
	idUser1 := h.ctx.Value(keyId).(int)
	log.Println("id = ", idUser1)
	// получение email пользователя, с которым создаем чат, из контекста
	emailUser2 := h.ctx.Value(keyEmail).(string)

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
		err = h.service.WriteInChat(msg, idUser1, emailUser2)
		if err != nil {
			log.Fatalln("error: сообщение не сохранено: ", err)
		}

		// отправляем сообщение
		go h.sendMessage(msg, keyClients)
	}
}
