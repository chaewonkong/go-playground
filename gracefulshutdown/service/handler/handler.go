package handler

import (
	"context"
	"fmt"
	"gracefulshutdown/pkg/redis"
	"log"
	"net/http"
)

type Handler struct {
	RedisClient *redis.RedisClient
}

func New(rc *redis.RedisClient) *Handler {
	return &Handler{
		RedisClient: rc,
	}
}

func (h *Handler) Incr(w http.ResponseWriter, r *http.Request) {
	key := "reqkey"

	log.Printf("Request key %s", key)

	ctx := r.Context()

	// mock latency
	// redis with 10 seconds delay
	result, err := h.RedisClient.Incr(ctx, key).Result()

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
	result, err := h.RedisClient.Get(ctx, key).Result()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Decremented count for %s:%v", key, result)
	_, _ = w.Write([]byte(fmt.Sprintf("Decremented count for %s:%v", key, result)))
}
