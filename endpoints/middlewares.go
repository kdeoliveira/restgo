package endpoints

import "net/http"

type Cors struct {
	Remote string
}

func (cors *Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		response.Header().Set("Content-Type", "application/json")
		response.Header().Set("Access-Control-Allow-Origin", cors.Remote)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(response, request)
	})
}
