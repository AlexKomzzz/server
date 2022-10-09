package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	chat "github.com/AlexKomzzz/server"
)

type HistoryResp struct {
	History []chat.Message
}

// создание чата с другим пользователем по его email
func (h *Handler) getChat(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "invalid method: no GET", http.StatusBadRequest)
			return
		}

		//// для проверки без идентификации
		// h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, 1)
		////

		// var historyChat []chat.Message

		// // выделим email из url
		// // получим мапу из параметров указанных в url с помощью "?"
		set, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
			return
		}

		// log.Println("set = ", set)

		/////// ДЛЯ ТЕСТА
		// set["email"] = append(set["email"], "bobik")
		// ///////

		var emailUser2 string
		// // из мапы вытащим значение email
		if _, ok := set["email"]; ok {
			emailUser2 = set["email"][0]
		}
		// log.Println("emailUser2 = ", emailUser2)

		// u, _ := url.Parse(r.URL)
		// emailUser2 := r.URL.Fragment
		// log.Println("emailUser2 = ", emailUser2)
		// записать в контекст
		h.webClient.ctx = context.WithValue(h.webClient.ctx, keyEmail, emailUser2)

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
}

// открываем соединение, в цикле читаем сообщения и парсим в структуру
func (h *Handler) ChatTwoUser(w http.ResponseWriter, r *http.Request) {

	////////
	//upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	////////

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("error connect: ", err)
	}
	defer conn.Close()
	log.Println("Client connected:", conn.RemoteAddr().String())

	//clnt.ctx = context.WithValue(clnt.ctx, keyName, "Alex")
	//h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, 1)

	// вытащим username пользователя из контекста
	//username := h.webClient.ctx.Value(keyName).(string)

	//userId := h.webClient.ctx.Value(keyId).(int)
	//username, err := h.service.GetUsername(userId)
	username, err := "Alex", nil
	if err != nil {
		log.Fatalln("error: не получен username по id: ", err)
	}

	// сохраняем соединение
	h.webClient.clients[conn] = true
	defer delete(h.webClient.clients, conn)

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
		//h.service.WriteInChat(msg, )

		go h.MyWriteMessage(msg) // отправляем сообщение
		// go messageHandler(message) // выводим сообщение

		//////////////////////////////messageHandler(msg.Message)

	}
}
