package main

import (
	"context"
	"fmt"
	"gracefulshutdown/service"
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	time.Sleep(1 * time.Minute)
	return rc.client.Incr(ctx, key)
}

func (rc *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return rc.client.Get(ctx, key)
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := "reqkey"

	log.Printf("Request key %s", key)

	ctx := context.Background()

	// mock latency
	// redis with 10 seconds delay
	result, err := h.client.Get(ctx, key).Result()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Decremented count for %s:%v", key, result)
	_, _ = w.Write([]byte(fmt.Sprintf("Decremented count for %s:%v", key, result)))
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rc := &RedisClient{
		client: client,
	}
	h := &Handler{
		client: rc,
	}

	cancellableMux := NewDefaultMux()
	nonCancellableMux := NewDefaultMux()
	cancellableMux.HandleFunc("/read", h.Get)     // no latency
	nonCancellableMux.HandleFunc("/incr", h.Incr) // latency 10 seconds

	cancellableSvr := &http.Server{
		Addr:    ":8080",
		Handler: cancellableMux,
	}

	nonCancellableSvr := &http.Server{
		Addr:    ":8081",
		Handler: nonCancellableMux,
	}

	cancellableTask := service.NewTask("cancellable-http-server", true, cancellableSvr)
	nonCancellableTask := service.NewTask("non-cancellable-http-server", false, nonCancellableSvr)

	svc := &service.Service{
		Tasks: []service.Task{cancellableTask, nonCancellableTask},
	}

	ctx := context.Background()
	svc.Init(ctx)

	go func() {
		log.Println("Starting the service...")
		if err := svc.Run(); err != nil {
			log.Println(err)
		}

	}()

	// cancel the service
	<-c
	// 10초 이내에 종료되지 않으면 강제 종료
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	svc.Close(ctx)
}

// 기대 결과
// 1. cancellable: 즉시 종료되어야 함 ctx의 cancelFunc 호출.
// 2. non-cancellable: 10초 이내에 종료되지 않으면 강제 종료.
//      Graceful Shutdown으로 10초 내에는 작업 종료를 대기.
