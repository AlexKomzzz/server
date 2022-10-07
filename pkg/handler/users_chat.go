package handler

import (
	"net/http"

	chat "github.com/AlexKomzzz/server"
)

// создание чата с другим пользователем по его email
func (h *Handler) getChat(w http.ResponseWriter, r *http.Request) {

	// проверка метода
	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// выделим email из url
	//r.URL.Parse()
	emailUser2 := "asd"

	// вытащим id пользователя из контекста
	idUser := h.webClient.ctx.Value(keyId).(int)

	//
	historyChat := make([]*chat.Message, 0)
	var err error

	// получение истории чата с пользователем
	historyChat, err = h.service.GetChat(historyChat, idUser, emailUser2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
