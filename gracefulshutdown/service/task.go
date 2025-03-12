package service

import (
	"context"
	"fmt"
	"net/http"
)

type Task interface {
	Run(ctx context.Context, sigTerm <-chan struct{}) error
	Cancellable() bool
}

type TaskImpl struct {
	name        string
	cancellable bool
	svr         *http.Server // TODO: http.Server 뿐만 아니라 grpc.Server도 받을 수 있도록 수정
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
func (t *TaskImpl) Run(ctx context.Context, sigTerm <-chan struct{}) error {
	// Run server
	errCh := make(chan error, 1)
	go func() {
		err := t.svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error, %s: %w", t.name, err)
		}
	}()

	// Cancellation
	select {
	case <-ctx.Done(): // Graceful shutdown 시도
		if err := t.svr.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown error, %s: %w", t.name, err)
		}
		return nil
	case <-sigTerm:
		return fmt.Errorf("force shutdown error, %s", t.name)
	case err := <-errCh:
		return err
	}
}

func (t *TaskImpl) Cancellable() bool {
	return t.cancellable
}
