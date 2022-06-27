package handler

import (
	"net/http"
)

type ControllerMethod func(http.ResponseWriter, *http.Request)

type Methods uint8

const (
	GET Methods = iota
	POST
	DELETE
)

func (m Methods) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case DELETE:
		return "DELETE"
	default:
		return "GET"
	}
}
