package service

import (
	"context"
	"fmt"
	"net/http"
)

type Task interface {
	Run(ctx context.Context) error
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
func (t *TaskImpl) Run(ctx context.Context) error {
	t.ctx = ctx

	errCh := make(chan error, 1)
	// Run server
	go func() {
		if err := t.svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error, %s: %w", t.name, err)
		}
	}()

	select {
	case <-t.ctx.Done(): // cancellable이고 cancel된 경우 서버 종료
		return t.svr.Shutdown(t.ctx)
	case err := <-errCh: // 에러 발생인 경우 에러 반환하고 서버 종료
		return err
	}
}

func (t *TaskImpl) GracefulShutdown(ctx context.Context) error {
	return t.svr.Shutdown(ctx)
}

func (t *TaskImpl) Cancellable() bool {
	return t.cancellable
}

func (t *TaskImpl) String() string {
	return t.name
}
