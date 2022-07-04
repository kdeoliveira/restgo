package main

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/kdeoliveira/restgo/endpoints"
	"github.com/kdeoliveira/restgo/server"
	ws2 "github.com/kdeoliveira/restgo/ws"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	port   int
	addr   string
	remote string
)

func init() {
	var err error
	port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 8000
	}
	addr = os.Getenv("API_ADDR")
	remote = os.Getenv("API_REMOTE")
}

func main() {

	ws := ws2.WsServer{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			WriteBufferSize: 1024,
		},
	}

	httpServer := server.New(addr, port)

	httpServer.AddHandler(endpoints.GET, "/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := ws.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		go endpoints.WsSendDateTime(conn)
	})

	httpServer.AddHandler(endpoints.GET, "/", endpoints.HomeHandler)

	httpServer.AddHandler(endpoints.GET, "/time", endpoints.CurrentTime)

	httpServer.AddHandler(endpoints.GET, "/locale", endpoints.AllLocales)

	cors := endpoints.Cors{Remote: remote}

	httpServer.AddMiddleware(cors.Middleware)

	sv, done := httpServer.Serve(func(server *http.Server) {
		log.Println("Server started successfully")
	})

	var wait = time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	select {
	case <-done:
		log.Println("Attempting to stop gracefully")
		err := sv.Shutdown(ctx)
		if err != nil {
			log.Fatal("Error occurred! ", err)
		}
	}

	log.Println("Shutting down gracefully...")
	os.Exit(0)

}
