package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	chat "github.com/AlexKomzzz/server"
)

type HistoryResp struct {
	History []chat.Message
}

// создание чата с другим пользователем по его email
func (h *Handler) getChat(w http.ResponseWriter, r *http.Request) {

	// проверка метода
	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// выделим email из url
	// получим мапу из параметров указанных в url с помощью "?"
	set, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
		return
	}
	// из мапы вытащим значение email
	emailUser2 := set["email"][0]

	// вытащим id пользователя из контекста
	idUser := h.webClient.ctx.Value(keyId).(int)

	//
	historyChat := make([]chat.Message, 0)

	// получение истории чата с пользователем
	historyChat, err = h.service.GetChat(historyChat, idUser, emailUser2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// http ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&HistoryResp{
		History: historyChat,
	})
}
