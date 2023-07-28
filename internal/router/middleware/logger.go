package middleware

import (
	"log"
	"net/http"
)

func LogRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.Host, r.Method, r.URL.String())

		f(w, r)
	}
}
