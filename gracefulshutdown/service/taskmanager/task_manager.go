package taskmanager

import (
	"context"
	"errors"
	"gracefulshutdown/service"
	"sync"
)

var _ service.Application = (*TaskManager)(nil)

type TaskManager struct {
	wg    sync.WaitGroup
	errCh chan error
	tasks []func(context.Context) error
}

func (t *TaskManager) Run() error {
	ctx := context.Background()
	for _, task := range t.tasks {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			err := task(ctx)
			if err != nil {
				t.errCh <- err
			}
		}()
	}

	return nil
}

func (t *TaskManager) Shutdown(ctx context.Context) error {
	// 실행 중인 작업이 종료되는 것을 감지한다
	done := make(chan struct{})
	go func() {
		t.wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		// ctx가 종료되면 app을 종료한다
		return errors.New("app shutdown by ctx timeout")
	case <-done:
		// 모든 작업이 종료되면 app을 종료한다
		return nil
	}
}

func FindMatch(ctx context.Context) error {
	return nil
}
