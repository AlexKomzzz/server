package handler

import (
	//"context"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type myCtx string

// var keyName, keyId myCtx = "username", "userId"
var keyId, keyEmail myCtx = "userId", "email"

// поиск email в URL и проверка идентификации
func (h *Handler) parseEmailAndIdentity(next http.Handler) http.Handler {
	return h.parseEmail(h.userIdentity(next))
}

// проверка идентификации для Handler
// парсинг хедера, определение JWT, определение id
func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization") // выделяем из заголовка поле "Authorization"
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if userId < 1 {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		// запись idUser в контекст
		h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, userId)

		// username, err := h.service.GetUsername(userId)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }

		// // запись username в контекст
		// h.webClient.ctx = context.WithValue(h.webClient.ctx, keyName, username)
		//c.Set("userId", userId)

		next.ServeHTTP(w, r)
	})
}

// проверка идентификации для HandlerFunc
/*func (h *Handler) userIdentityHF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization") // выделяем из заголовка поле "Authorization"
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if userId < 1 {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		// запись idUser в контекст
		h.webClient.ctx = context.WithValue(h.webClient.ctx, keyId, userId)

		// username, err := h.service.GetUsername(userId)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }

		// // запись username в контекст
		// h.webClient.ctx = context.WithValue(h.webClient.ctx, keyName, username)
		//c.Set("userId", userId)

		next(w, r)
	}
}*/

// парсинг URL в поиках email
func (h *Handler) parseEmail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// выделим email из url
		// получим мапу из параметров указанных в url с помощью "?"
		set, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
			return
		}
		var emailUser2 string
		// // из мапы вытащим значение email
		if _, ok := set["email"]; ok {
			emailUser2 = set["email"][0]
		}
		if emailUser2 == "" {
			log.Println("Значение email не передано в URL")
			//w.Write([]byte("error: не передано значние email в URL"))

		}

		h.webClient.ctx = context.WithValue(h.webClient.ctx, keyEmail, emailUser2)

		next.ServeHTTP(w, r)
	})
}
