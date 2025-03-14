package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Task interface {
	Run(ctx context.Context, termCh <-chan struct{}) error
	GracefulShutdown(ctx context.Context) error
	Cancellable() bool
	String() string
}

type TaskImpl struct {
	name        string
	cancellable bool
	ctx         context.Context
	svr         *http.Server
}

func NewTask(name string, cancellable bool, svr *http.Server) Task {
	return &TaskImpl{
		name:        name,
		cancellable: cancellable,
		svr:         svr,
	}
}

// Run runs the task
//   - ctx: cancellable 또는 non-cancellable context
func (t *TaskImpl) Run(ctx context.Context, termCh <-chan struct{}) error {
	t.ctx = ctx

	errCh := make(chan error, 1)
	// Run server
	go func() {
		errCh <- t.svr.ListenAndServe()
	}()

	select {
	case <-t.ctx.Done(): // cancellable이고 cancel된 경우 서버 종료
		log.Printf("ctx done: %s", t.name)
		return t.svr.Shutdown(t.ctx)
	case err := <-errCh: // 에러 발생인 경우 에러 반환하고 서버 종료
		log.Printf("err: %s", t.name)

		return err
	case <-termCh: // 종료 신호가 왔을 경우 서버 종료: 비정상 종료에 해당
		return fmt.Errorf("server shutdown error, %s: %w", t.name, t.svr.Close())
	}

	// errCh := make(chan error, 1)
	// go func() {
	// 	err := t.svr.ListenAndServe()
	// 	if err != nil && err != http.ErrServerClosed {
	// 		errCh <- fmt.Errorf("server error, %s: %w", t.name, err)
	// 	}
	// }()

	// // ctx.Done() 이면 바로 종료

	// // Cancellation
	// select {
	// case <-ctx.Done(): // Graceful shutdown 시도
	// 	if err := t.svr.Close(); err != nil {
	// 		return fmt.Errorf("graceful shutdown error, %s: %w", t.name, err)
	// 	}
	// 	return nil
	// case <-sigTerm:
	// 	err := t.svr.Close() // 서버 즉시 종료
	// 	if err != nil {
	// 		return fmt.Errorf("close error, %s: %w", t.name, err)
	// 	}
	// 	return nil
	// case err := <-errCh:
	// 	return err
	// }
}

func (t *TaskImpl) GracefulShutdown(ctx context.Context) error {
	// // cancellable은 t.ctx의 cancel로 종료
	// if t.cancellable {
	// 	return nil
	// }

	// // non-cancellable은 graceful shutdown
	err := t.svr.Shutdown(ctx)
	return err
}

func (t *TaskImpl) Cancellable() bool {
	return t.cancellable
}

func (t *TaskImpl) String() string {
	return t.name
}
