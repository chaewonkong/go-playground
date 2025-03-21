package service

import (
	"context"
	"fmt"
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
	app         App

	shutdownFunc func(ctx context.Context) error
}

func NewTask(name string, cancellable bool, app App) Task {
	return &TaskImpl{
		name:        name,
		cancellable: cancellable,
		app:         app,
	}
}

// Run runs the task
//   - ctx: cancellable 또는 non-cancellable context
func (t *TaskImpl) Run(ctx context.Context) error {
	t.ctx = ctx

	errCh := make(chan error, 1)
	// Run server
	go func() {
		if err := t.app.Run(); err != nil {
			errCh <- fmt.Errorf("server error, %s: %w", t.name, err)
		}
	}()

	select {
	case <-t.ctx.Done(): // cancellable이고 cancel된 경우 서버 종료
		return t.app.Stop()
	case err := <-errCh: // 에러 발생인 경우 에러 반환하고 서버 종료
		return err
	}
}

func (t *TaskImpl) GracefulShutdown(ctx context.Context) error {
	return t.app.GracefulStop(ctx)
}

func (t *TaskImpl) Cancellable() bool {
	return t.cancellable
}

func (t *TaskImpl) String() string {
	return t.name
}

type App interface {
	Run() error
	Stop() error
	GracefulStop(ctx context.Context) error
}
