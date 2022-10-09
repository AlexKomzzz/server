package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func (h *Handler) test(w http.ResponseWriter, r *http.Request) {

	// выделим email из url
	// получим мапу из параметров указанных в url с помощью "?"
	set, _ := url.ParseQuery(r.URL.RawQuery)

	// log.Println("set = ", set)

	// ///// ДЛЯ ТЕСТА
	// // set["email"] = append(set["email"], "bobik")
	// ///////

	var emailUser2 string
	// // из мапы вытащим значение email
	if _, ok := set["email"]; ok {
		emailUser2 = set["email"][0]
	}

	//emailUser2 := r.URL.Fragment
	// log.Println("emailUser2 = ", emailUser2)

	// вытащим id пользователя из контекста
	// idUser := h.webClient.ctx.Value(keyId).(int)

	//получение истории чата с пользователем
	// historyChat, err = h.service.GetChat(idUser, emailUser2)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// http ответ
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(&HistoryResp{
	// 	History: historyChat,
	// })

	json.NewEncoder(w).Encode(&emailUser2)

}
