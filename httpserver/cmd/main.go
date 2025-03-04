package main

import (
	"httpserver/pkg/httpserver"
	"log"
	"net/http"
	"time"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		log.Println("Request processed in", time.Since(now), r.Method, r.URL.Path)
	})
}

func healthChecker(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func wildcardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("Wildcard handler, id: %s\n", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Wildcard"))
}

func main() {
	// init svrConfig
	port := "8080"

	// init server
	app := httpserver.New()
	app.Get("/", healthChecker)
	app.Get("/users/{id}", wildcardHandler)

	// add middlewares
	app.Use(logMiddleware)

	// run
	log.Fatal(app.Listen(port))
}
