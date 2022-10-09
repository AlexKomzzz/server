package handler

import (
	"log"
	"net/http"
	"time"

	chat "github.com/AlexKomzzz/server"
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
	username, err := "Alex", nil
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}
	////////////////////////

	// сохраняем соединение
	h.webClient.clients[conn] = true
	defer delete(h.webClient.clients, conn)

	// получение id текущего пользователя из контекста
	idUser1 := h.webClient.ctx.Value(keyId).(int)
	// получение email пользователя, с которым создаем чат, из контекста
	emailUser2 := h.webClient.ctx.Value(keyEmail).(string)

	// получение истории чата из БД
	historyChat, err := h.service.GetChat(idUser1, emailUser2)
	if err != nil {
		log.Fatalln("error: не получена история чата: ", err)
	}

	// передача истории клиентам
	for _, msg := range historyChat {
		go h.MyWriteMessage(&msg) // возможна блокировка
	}

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
		h.service.WriteInChat(msg, idUser1, emailUser2)

		go h.MyWriteMessage(msg) // отправляем сообщение

	}
}
