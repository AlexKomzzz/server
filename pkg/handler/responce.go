package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) newErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err string) {

	log.Printf("error: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(&errorResponse{
		Message: err,
	})

}
