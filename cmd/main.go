package main

import (
	"context"
	"github.com/kdeoliveira/restgo/controller"
	"github.com/kdeoliveira/restgo/server"
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

	httpServer := server.New(addr, port)

	httpServer.AddHandler(controller.GET, "/", controller.HomeHandler)

	httpServer.AddHandler(controller.GET, "/time", controller.CurrentTime)

	cors := controller.Cors{Remote: remote}

	httpServer.AddMiddleware(cors.Middleware)

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
