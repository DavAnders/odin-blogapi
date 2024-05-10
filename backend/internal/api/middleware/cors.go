package middleware

import (
	"log"
	"net/http"
)

// CORS middleware
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // Allow requests from any origin
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
			log.Println("Pre-flight request received")
            w.WriteHeader(http.StatusOK)
            return
        }

		log.Printf("Request received: %s %s", r.Method, r.URL.Path) // Log the method and path of each request

        next.ServeHTTP(w, r)
    })
}
