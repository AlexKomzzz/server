package lol

import (
	"github.com/gorilla/websocket"
)

type Lol struct {
	clients map[*websocket.Conn]bool
}

func NewLol(clients map[*websocket.Conn]bool) *Lol {
	return &Lol{clients: clients}
}
