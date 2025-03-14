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
	// wg               sync.WaitGroup
}

func (s *Service) Init(ctx context.Context) {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.unCancellableCtx = context.WithoutCancel(ctx)
	s.closeCh = make(chan struct{})
}

func (s *Service) Run() error {
	errCh := make(chan error, len(s.Tasks))

	var wg sync.WaitGroup

	for _, task := range s.Tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			if task.Cancellable() {
				err := task.Run(s.ctx, s.closeCh)
				if err != nil {
					errCh <- fmt.Errorf("run %s: %w", task.(fmt.Stringer).String(), err)
				}
			} else {
				err := task.Run(s.unCancellableCtx, s.closeCh)
				if err != nil {
					errCh <- fmt.Errorf("run %s: %w", task.(fmt.Stringer).String(), err)
				}
			}
		}(task)
	}

	// // s.wg.Wait()
	// <-s.closeCh // 종료하라는 신호가 왔다
	// s.cancel()  // 종료 가능한 context는 종료
	// for _, task := range s.Tasks {
	// 	if !task.Cancellable() {
	// 		go func(task Task) {
	// 			err := task.Stop(s.unCancellableCtx)
	// 			if err != nil {
	// 				errCh <- fmt.Errorf("%s: %w", task.(fmt.Stringer).String(), err)
	// 			}
	// 		}(task)
	// 	}
	// }
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
				err := task.GracefulShutdown(s.unCancellableCtx)
				if err != nil {
					errCh <- fmt.Errorf("%s: %w", task.(fmt.Stringer).String(), err)
				}
			}(task)
		}
	}

	go func() {
		wg.Wait()
		close(doneCh)
		// log.Println("Gracefully shutdown")
	}()

	select {
	case <-doneCh:
		// 정상 종료 완료
		log.Println("Gracefully shutdown")
	case <-termCtx.Done():
		// 강제 종료
		close(s.closeCh)
		log.Println("Force shutdown")
	}

	close(errCh)

	for err := range errCh {
		log.Println(err)
	}

	// 1. 모든 서비스가 다 종료된 경우
	// 즉시 return
	// 2. 서비스가 다 종료되지 않은 경우
	// 강제 종료 대기

	// 강제 종료의 순간이 왔다
	// <-termCtx.Done()
	// close(s.closeCh)

	// for err := range errCh {
	// 	log.Println(err)
	// }

	// select {
	// case <-s.closeCh:
	// 	errCh := make(chan error, len(s.Tasks))

	// 	// 모든 작업이 종료되는 것 대기
	// 	s.cancel() // 종료 가능한 context는 종료

	// 	// 종료되지 않은 작업에 대해 강제 종료
	// 	for _, task := range s.Tasks {
	// 		if !task.Cancellable() {
	// 			go func(task Task) {
	// 				err := task.Stop(s.unCancellableCtx)
	// 				if err != nil {
	// 					errCh <- fmt.Errorf("%s: %w", task.(fmt.Stringer).String(), err)
	// 				}
	// 			}(task)
	// 		}
	// 	}
	// 	for err := range errCh {
	// 		log.Println(err)
	// 	}
	// 	log.Println("Gracefully shutdown")
	// case <-termCtx.Done():
	// 	// 강제 종료
	// 	close(s.closeCh)
	// 	log.Println("Force shutdown")
	// }
}
