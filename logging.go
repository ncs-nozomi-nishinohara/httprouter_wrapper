package httprouter_wrapper

import (
	"log"
	"net/http"
)

// アクセスログ
func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		h.ServeHTTP(w, r)
	})
}
