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

func main() {
	// init svrConfig

	// init serveMux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", healthChecker)

	// init server
	svr := httpserver.New(mux, &http.Server{Addr: ":8080"})

	// add middlewares
	svr.Use(authMiddleware, logMiddleware)

	// run
	log.Fatal(svr.Run())
}
