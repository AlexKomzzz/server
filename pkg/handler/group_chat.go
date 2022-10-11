package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	idUser := h.webClient.ctx.Value(keyId).(int)

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

	idGroup, err := h.service.CreateGroup(title_group, idUser)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// отправим клиенту id группового чата
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\n\t\"idGroup\": \"%d\"\n}", idGroup)))
}
