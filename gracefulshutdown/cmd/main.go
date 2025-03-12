package main

import (
	"context"
	"gracefulshutdown/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// simulate a long request
		time.Sleep(10 * time.Minute)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	cancellableSvr := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	nonCancellableSvr := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	cancellableTask := service.NewTask("cancellable-http-server", true, cancellableSvr)
	nonCancellableTask := service.NewTask("non-cancellable-http-server", false, nonCancellableSvr)

	svc := &service.Service{
		Tasks: []service.Task{cancellableTask, nonCancellableTask},
	}

	ctx := context.Background()
	svc.Init(ctx)

	log.Println("Starting the service...")
	if err := svc.Run(); err != nil {
		log.Println(err)
	}

	// cancel the service
	<-c
	// 10초 이내에 종료되지 않으면 강제 종료
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	svc.Close(ctx)
}
