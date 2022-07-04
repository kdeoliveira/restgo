package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WsServer struct {
	Upgrader websocket.Upgrader
	Handlers http.HandlerFunc
}
