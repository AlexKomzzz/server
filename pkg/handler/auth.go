package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	chat "github.com/AlexKomzzz/server"
)

// Обработчик для регистрации пользователя
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user chat.User

	// парсим тело запроса в структуру пользователя
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid input body", http.StatusBadRequest)
		return
	}

	// по данным пользователя заносим в БД и получаем id
	id, err := h.service.CreateUser(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("error createUser handler: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("\"id\": \"%d\"", id)))
	// c.JSON(http.StatusOK, gin.H{
	// 	"id": id,
	// })

}

type InUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResp struct {
	Token string `json:"token"`
}

// Обработчик для аутентификации пользователя
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var user InUser

	// парсим тело запроса в структуру пользователя
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid input body", http.StatusBadRequest)
		return
	}

	// по данным пользователя заносим в БД и получаем id
	token, err := h.service.GenerateToken(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(fmt.Sprintf("{\n\t\"token\": \"%s\"\n}", token)))
	json.NewEncoder(w).Encode(&TokenResp{
		Token: token,
	})

	// c.JSON(http.StatusOK, gin.H{
	// 	"token": token,
	// })
}

// func (h *Handler) In(w http.ResponseWriter, r *http.Request) {

// 	w.Write([]byte(fmt.Sprintf("\"token\"")))
// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"token": token,
// 	// })
// }
