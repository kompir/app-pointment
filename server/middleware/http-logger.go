package middleware

import (
	"log"
	"net/http"
	"strings"
)

func HTTPLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", strings.ToUpper(r.Method), r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
