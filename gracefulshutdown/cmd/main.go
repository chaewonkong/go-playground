package main

import (
	"context"
	"gracefulshutdown/pkg/redis"
	"gracefulshutdown/service"
	"gracefulshutdown/service/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewDefaultMux(funcs ...http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	return mux
}

func NewServer(port string, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

type Server struct {
	*http.Server
}

func (s *Server) Run() error {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

func (s *Server) Stop() error {
	return s.Close()
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	rc := redis.NewClient(1 * time.Minute)
	h := handler.New(rc)

	cancellableMux := NewDefaultMux()
	nonCancellableMux := NewDefaultMux()
	cancellableMux.HandleFunc("/read", h.Get)     // no latency
	nonCancellableMux.HandleFunc("/incr", h.Incr) // latency 10 seconds

	cancellableSvr := NewServer("8080", cancellableMux)
	nonCancellableSvr := NewServer("8081", nonCancellableMux)

	cancellableTask := service.NewTask("cancellable-http-server", true, &Server{cancellableSvr})
	_ = cancellableTask
	nonCancellableTask := service.NewTask("non-cancellable-http-server", false, &Server{nonCancellableSvr})

	// Create a new service

	svc := &service.Service{
		Tasks: []service.Task{nonCancellableTask},
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

/* 테스트 케이스 정리
1. 즉시 종료 ctx는 inturrupt 발생시 즉시 종료되는가? ✅
2. 대기context는 10초 이내에 종료되지 않으면 강제 종료되는가? 이 경우 로그가 남는가?
3. 대기 context는 10초 이내에 만약 graceful shutdown이 완료되면 종료되는가? 이 경우 redis 저장이 완료되는가? ✅
*/
