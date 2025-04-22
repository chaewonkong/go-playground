package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const tickDuration = 1 * time.Second

func main() {
	err := run()
	if err != nil {
		slog.Default().Error("failed to run", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	rdb := redis.NewClient(&redis.Options{})
	defer func() { _ = rdb.Close() }()

	// redis 연결 확인
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	// logger
	logger := slog.Default()

	// signal을 통해 종료 신호를 받는다
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// queue에 등록된 티켓 조회
				ticketIDs, err := rdb.LRange(context.Background(), "ticket", 0, -1).Result()
				if err != nil {
					logger.Error("failed to get ticket IDs from redis: " + err.Error())
				}

				workers := []Worker{}
				for _, ticket := range ticketIDs {
					// 실제로는 일정 갯수로 분배
					workers = append(workers, Worker{tickets: []string{ticket}})
				}

				for _, worker := range workers {
					wg.Add(1)
					go func(w Worker) {
						defer wg.Done()
						// job을 worker가 실행
						w.FindMatches()
					}(worker)
				}
			}
		}
	}(ctx)

	// 종료 신호 발생
	<-ctx.Done()

	// graceful shutdown에 대한 timeout을 설정
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 진행 중인 작업들의 종료를 기다린다. (context timeout이 되기 전까지)
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done(): // ctx가 종료되면 app을 비정상 종료한다
		return errors.New("app shutdown by ctx timeout")
	case <-doneCh: // 모든 작업이 종료되면 app을 정상 종료한다
		return nil
	}
}

type Worker struct {
	tickets []string
}

func (w *Worker) FindMatches() {
	// match 찾는 로직
	// redis에 발행하는 로직

	return
}
