package handler

import "net/http"

// создание чата с другим пользователем по его email
func (h *Handler) createChat(w http.ResponseWriter, r *http.Request) {

	// проверка метода
	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// выделим email из url
	//r.URL.Parse()

	// получить id второго клиента по его email
}
