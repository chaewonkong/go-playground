package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Service struct {
	Tasks            []Task
	ctx              context.Context
	unCancellableCtx context.Context
	cancel           context.CancelFunc
}

func (s *Service) Init(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.unCancellableCtx = context.WithoutCancel(ctx)
	// s.closeCh = make(chan struct{})
}

func (s *Service) Run() error {
	errCh := make(chan error, len(s.Tasks))

	var wg sync.WaitGroup

	for _, task := range s.Tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			if task.Cancellable() {
				err := task.Run(s.ctx)
				if err != nil {
					errCh <- fmt.Errorf("run %s: %w", task.(fmt.Stringer).String(), err)
				}
			} else {
				err := task.Run(s.unCancellableCtx)
				if err != nil {
					errCh <- fmt.Errorf("run %s: %w", task.(fmt.Stringer).String(), err)
				}
			}
		}(task)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		joinedErr := errors.Join(err)
		if joinedErr != nil {
			return joinedErr
		}
	}

	return nil
}

// Close closes the service
//   - termCtx에는 비관적 종료를 강제하는 timeout이 설정되어 있음 가정
func (s *Service) Close(termCtx context.Context) {
	var wg sync.WaitGroup

	// 종료 가능한 context는 종료
	s.cancel()

	errCh := make(chan error, len(s.Tasks))
	doneCh := make(chan struct{}) // 모든 작업이 종료되었는지 확인

	// 종료되지 않은 작업에 대해 강제 종료
	for _, task := range s.Tasks {
		if !task.Cancellable() {
			wg.Add(1)
			go func(task Task) {
				defer wg.Done()
				err := task.GracefulShutdown(termCtx)

				// 비정상 종료 등
				if err != nil {
					errCh <- fmt.Errorf("error to shutdown server: %s, error: %w", task.String(), err)
				}
			}(task)
		}
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// 정상 종료 완료
		log.Println("gracefully shutdown")
	case <-termCtx.Done():
		// 강제 종료
		<-doneCh
		log.Println("force shutdown")
	}

	close(errCh)

	for err := range errCh {
		log.Println(err)
	}
}
