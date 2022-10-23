package handler

import (
	//"context"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type myCtx string

// keyId, keyIdUser2, keyIdGroup, keyIdChat, keyTitle myCtx = "userId", "idUser2", "idGroup", "idChat", "title"
var keyId, keyIdUser2, keyIdGroup, keyIdChat, keyTitle myCtx = "userId", "idUser2", "idGroup", "idChat", "title"

// поиск email в URL и проверка идентификации
func (h *Handler) identityAndParseURL(next http.Handler) http.Handler {
	return h.userIdentity(h.parseURL(next))
}
func (h *Handler) identityAndParseURLHF(next http.HandlerFunc) http.HandlerFunc {
	return h.userIdentityHF(h.parseURLHF(next))
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
		h.ctx = context.WithValue(h.ctx, keyId, userId)

		next.ServeHTTP(w, r)
	})
}

// проверка идентификации для HandlerFunc
func (h *Handler) userIdentityHF(next http.HandlerFunc) http.HandlerFunc {
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

		log.Printf("idUser = %d", userId)
		// запись idUser в контекст
		h.ctx = context.WithValue(h.ctx, keyId, userId)

		next(w, r)
	}
}

// парсинг URL в поиках email
/*func (h *Handler) parseEmail(next http.Handler) http.Handler {
	return h.parseEmailHF(next.ServeHTTP)
}

// парсинг URL в поиках email для HandlerFunc
func (h *Handler) parseEmailHF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		h.ctx = context.WithValue(h.ctx, keyEmail, emailUser2)

		next(w, r)
	}
}*/

// парсинг URL для Handler
func (h *Handler) parseURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// выделим idGroup из url
		// получим мапу из параметров указанных в url с помощью "?"
		set, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
			return
		}

		// если в URL передан фрагмент idGroup запишем его в контекст
		if _, ok := set["idGroup"]; ok {

			idGroupStr := set["idGroup"][0]

			// конвертация idGroup из стороковго типа в целочисленный
			idGroup, err := strconv.Atoi(idGroupStr)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idGroup: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdGroup, idGroup)

		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdGroup, nil)
		}

		// если в URL передан фрагмент email запишем его в контекст
		if _, ok := set["idUser2"]; ok {

			idUser2Str := set["idUser2"][0]

			// конвертация idUser2 из стороковго типа в целочисленный
			idUser2, err := strconv.Atoi(idUser2Str)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idUser2 from URL: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdUser2, idUser2)

			log.Printf("idUser2 определен = %d\n", idUser2)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdUser2, nil)
		}

		// если в URL передан фрагмент idChat запишем его в контекст
		if _, ok := set["idChat"]; ok {

			idChatStr := set["idChat"][0]

			// конвертация idChat из стороковго типа в целочисленный
			idChat, err := strconv.Atoi(idChatStr)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idChat from URL: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdChat, idChat)

			log.Printf("idChat определен = %d\n", idChat)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdChat, nil)
		}

		// если в URL передан фрагмент title запишем его в контекст
		if _, ok := set["title"]; ok {

			title := set["title"][0]

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyTitle, title)

			log.Printf("title определен = %s\n", title)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyTitle, nil)
		}

		next.ServeHTTP(w, r)
	})
}

// парсинг URL для HandkerFunc
func (h *Handler) parseURLHF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// выделим idGroup из url
		// получим мапу из параметров указанных в url с помощью "?"
		set, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, fmt.Errorf("error: invalid URL: %s", err).Error(), http.StatusBadRequest)
			return
		}

		// если в URL передан фрагмент idGroup запишем его в контекст
		if _, ok := set["idGroup"]; ok {

			idGroupStr := set["idGroup"][0]

			// конвертация idGroup из стороковго типа в целочисленный
			idGroup, err := strconv.Atoi(idGroupStr)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idGroup: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdGroup, idGroup)

		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdGroup, nil)
		}

		// если в URL передан фрагмент email запишем его в контекст
		if _, ok := set["idUser2"]; ok {

			idUser2Str := set["idUser2"][0]

			// конвертация idUser2 из стороковго типа в целочисленный
			idUser2, err := strconv.Atoi(idUser2Str)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idUser2 from URL: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdUser2, idUser2)

			log.Printf("idUser2 определен = %d\n", idUser2)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdUser2, nil)
		}

		// если в URL передан фрагмент idChat запишем его в контекст
		if _, ok := set["idChat"]; ok {

			idChatStr := set["idChat"][0]

			// конвертация idChat из стороковго типа в целочисленный
			idChat, err := strconv.Atoi(idChatStr)
			if err != nil {
				http.Error(w, fmt.Errorf("error: invalid idChat from URL: %s", err).Error(), http.StatusBadRequest)
				return
			}

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyIdChat, idChat)

			log.Printf("idChat определен = %d\n", idChat)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyIdChat, nil)
		}

		// если в URL передан фрагмент title запишем его в контекст
		if _, ok := set["title"]; ok {

			title := set["title"][0]

			// запись idGroup в контекст
			h.ctx = context.WithValue(h.ctx, keyTitle, title)

			log.Printf("title определен = %s\n", title)
		} else {
			// сбросим значение в контексте
			h.ctx = context.WithValue(h.ctx, keyTitle, nil)
		}

		next(w, r)
	}
}

// сравнение idUser1 и idUser2
func (h *Handler) compareIdUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// получение id текущего пользователя из контекста
		idUser1 := h.ctx.Value(keyId).(int)

		// получение idUser2 пользователя, с которым создаем чат, из контекста
		idUser2 := h.ctx.Value(keyIdUser2).(int)

		if idUser1 == idUser2 {
			h.newErrorResponse(w, r, http.StatusBadRequest, "id пользователя, переданный в URL, совпадает с собственным id")
		}
		next.ServeHTTP(w, r)
	})
}
