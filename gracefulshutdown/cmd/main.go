package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	return mux
}

type Handler struct {
	client *RedisClient
}

type RedisClient struct {
	client *redis.Client
}

func (rc *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	// mock latency
	time.Sleep(10 * time.Second)
	return rc.client.Incr(ctx, key)
}

func (rc *RedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	return rc.client.Decr(ctx, key)
}

func (h *Handler) Incr(w http.ResponseWriter, r *http.Request) {
	key := "reqkey"

	log.Printf("Request key %s", key)

	ctx := context.Background()

	// mock latency
	// redis with 10 seconds delay
	result, err := h.client.Incr(ctx, key).Result()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Incremented count for %s:%v", key, result)
	_, _ = w.Write([]byte(fmt.Sprintf("Incremented count for %s:%v", key, result)))
}

func (h *Handler) Decr(w http.ResponseWriter, r *http.Request) {
	key := "reqkey"

	log.Printf("Request key %s", key)

	ctx := context.Background()

	// mock latency
	// redis with 10 seconds delay
	result, err := h.client.Decr(ctx, key).Result()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Decremented count for %s:%v", key, result)
	_, _ = w.Write([]byte(fmt.Sprintf("Decremented count for %s:%v", key, result)))
}

func main() {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, os.Interrupt)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rc := &RedisClient{
		client: client,
	}
	h := &Handler{
		client: rc,
	}

	mux := NewDefaultMux()
	mux.HandleFunc("/incr", h.Incr) // latency 10 seconds
	mux.HandleFunc("/decr", h.Decr) // no latency

	svr := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("Starting the server on %s", svr.Addr)
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-term
	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Println("Gracefully shutdown")
}

// func main() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
// 		_, _ = w.Write([]byte("OK"))
// 	})
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// simulate a long request
// 		time.Sleep(10 * time.Minute)
// 		_, _ = w.Write([]byte("Hello, World!"))
// 	})

// 	cancellableSvr := &http.Server{
// 		Addr:    ":8080",
// 		Handler: mux,
// 	}

// 	nonCancellableSvr := &http.Server{
// 		Addr:    ":8081",
// 		Handler: mux,
// 	}

// 	cancellableTask := service.NewTask("cancellable-http-server", true, cancellableSvr)
// 	nonCancellableTask := service.NewTask("non-cancellable-http-server", false, nonCancellableSvr)

// 	svc := &service.Service{
// 		Tasks: []service.Task{cancellableTask, nonCancellableTask},
// 	}

// 	ctx := context.Background()
// 	svc.Init(ctx)

// 	go func() {
// 		log.Println("Starting the service...")
// 		if err := svc.Run(); err != nil {
// 			log.Println(err)
// 		}

// 	}()

// 	// cancel the service
// 	<-c
// 	// 10초 이내에 종료되지 않으면 강제 종료
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	svc.Close(ctx)
// }
