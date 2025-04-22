package service

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/Netflix/go-env"
)

type Application interface {
	// Run 애플리케이션을 실행한다
	Run() error

	// Shutdown 실행 중인 작업이 모두 종료되거나 ctx.Done()가 호출되면 App을 종료한다
	Shutdown(ctx context.Context) error
}

type AppCreatorFunc[CT any] func(config *CT) (Application, error)

func RunApplication[CT any](createApp AppCreatorFunc[CT], config *CT) error {
	_, err := env.UnmarshalFromEnviron(config)
	if err != nil {
		return errors.New("failed to unmarshal config")
	}

	app, err := createApp(config)
	if err != nil {
		return errors.New("failed to create app")
	}

	// signal을 통해 종료 신호를 받는다
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	errCh := make(chan error, 1)

	// App을 실행한다
	go func() {
		err := app.Run()
		if err != nil {
			errCh <- err
		}
	}()

	// 종료 처리
	select {
	case err := <-errCh:
		return err // nillable
	case <-ctx.Done(): // signal이 발생하면 ctx.Done()이 호출된다
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := app.Shutdown(ctx)
		return err // nillable
	}
}
