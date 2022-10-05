package handler

import "net/http"

func (h *Handler) createChat(w http.ResponseWriter, r *http.Request) {

	// проверка метода
	if r.Method != "POST" {
		http.Error(w, "invalid method: no POST", http.StatusBadRequest)
		return
	}

	// выделим email из url
	//r.URL.Parse()
}
