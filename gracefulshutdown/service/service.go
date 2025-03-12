package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Service struct {
	closeCh          chan struct{}
	Tasks            []Task
	ctx              context.Context
	unCancellableCtx context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
}

func (s *Service) Init(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.unCancellableCtx = context.WithoutCancel(ctx)
	s.closeCh = make(chan struct{})
}

func (s *Service) Run() error {
	errCh := make(chan error, len(s.Tasks))
	for _, task := range s.Tasks {
		s.wg.Add(1)
		go func(task Task) {
			defer s.wg.Done()
			if task.Cancellable() {
				err := task.Run(s.ctx, s.closeCh)
				errCh <- fmt.Errorf("%s: %w", task.(fmt.Stringer).String(), err)
			} else {
				err := task.Run(s.unCancellableCtx, s.closeCh)
				errCh <- fmt.Errorf("%s: %w", task.(fmt.Stringer).String(), err)
			}
		}(task)
	}

	s.wg.Wait()
	close(s.closeCh)
	close(errCh)

	for err := range errCh {
		joinedErr := errors.Join(err)
		if joinedErr != nil {
			return joinedErr
		}
	}

	return nil
}

func (s *Service) Close(ctx context.Context) {
	// ctx에는 timeout이 설정되어 있음 가정
	s.cancel()

	select {
	case <-s.closeCh:
		// 모든 작업이 종료되는 것 대기
		log.Println("Gracefully shutdown")
	case <-ctx.Done():
		// 타임아웃 대기 후 강제 종료
		close(s.closeCh)
		log.Println("Force shutdown")
	}
}
