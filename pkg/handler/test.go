package handler

import "net/http"

func (h *Handler) test(w http.ResponseWriter, r *http.Request) {
	h.service.Test()
}
