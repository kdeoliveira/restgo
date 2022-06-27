package main

import (
	"context"
	"github.com/kdeoliveira/restgo/handler"
	"github.com/kdeoliveira/restgo/server"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {

}

func main() {

	httpServer := server.New("127.0.0.1", 8000)

	httpServer.AddHandler(handler.GET, "/", handler.HomeHandler)

	httpServer.AddHandler(handler.GET, "/time", handler.CurrentTime)

	sv, done := httpServer.Serve(func(server *http.Server) {
		log.Println("Server started successfully")
	})

	var wait = time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	select {
	case <-done:
		err := sv.Shutdown(ctx)
		if err != nil {
			log.Fatal("Error occurred! ", err)
		}

	}

	log.Println("Shutting down gracefully...")
	os.Exit(0)

}
