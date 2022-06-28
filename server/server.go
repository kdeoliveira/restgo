package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kdeoliveira/restgo/controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	_ADDR = "127.0.0.1"
	_PORT = 8000
)

type Server struct {
	handlers []controller.Method
	port     int
	addr     string
	router   *mux.Router
}

func New(addr string, port int) *Server {
	if port == 0 {
		port = _PORT
	}
	if addr == "" {
		addr = _ADDR
	}

	return &Server{
		handlers: make([]controller.Method, 0),
		port:     port,
		addr:     addr,
		router:   mux.NewRouter(),
	}
}

func (server *Server) AddHandler(method controller.Methods, path string, handler controller.Method) {
	server.handlers = append(server.handlers, handler)
	// Note from gorilla/mux: there must be an OPTIONS method matcher for the middleware to set the headers
	server.router.HandleFunc(path, handler).Methods(method.String(), http.MethodOptions)
}

func (server *Server) AddMiddleware(mid mux.MiddlewareFunc) {
	server.router.Use(mid)
}

func (server *Server) Serve(callback func(server *http.Server)) (*http.Server, <-chan os.Signal) {
	done := make(chan os.Signal, 1)

	server.router.Use(mux.CORSMethodMiddleware(server.router))
	//https://stackoverflow.com/questions/19659600/how-to-use-gorilla-mux-with-http-timeouthandler
	// ReadTimeout is a timing constraint on the client http request imposed by the server from the moment
	// of initial connection up to the time the entire request body has been read.
	// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]

	// WriteTimeout is a time limit imposed on client connecting to the server via http from the
	// time the server has completed reading the request header up to the time it has finished writing the response.
	// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]

	sv := &http.Server{
		Handler:      server.router,
		Addr:         fmt.Sprintf("%s:%d", server.addr, server.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {

		if err := sv.ListenAndServe(); err != nil {
			log.Println(err)
			done <- os.Interrupt
		}
	}()

	signal.Notify(done, os.Interrupt)

	log.Printf("Initiating server @ http://%s:%d", server.addr, server.port)

	callback(sv)

	return sv, done
}
