package app

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// 환경변수에 저장된 config 파일의 path
const envConfigFilePath = "CONFIG_PATH"

// DeferFunc 정리 함수
type DeferFunc func()

// RunApplication Application을 실행하고 Signal을 받아 종료 처리를 담당하는 함수
func RunApplication(serviceName string, appFactory Factory) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT) // SIGTERM: k8s가 pod를 stop할 때 발생

	// global logger
	cleanUpLogger := initGlobalLogger(serviceName)
	defer cleanUpLogger()

	//config
	config := appFactory.NewConfig()
	err := ReadConfig(config)
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	// Application 생성
	app, err := appFactory.NewApp(config)
	if err != nil {
		zap.L().Fatal("failed to create application", zap.Error(err))
	}

	// Run
	err = app.Run()
	if err != nil {
		zap.L().Fatal("application run failed", zap.Error(err))
	}

	// Stop
	<-c
	err = app.Stop()
	if err != nil {
		zap.L().Fatal("application stop failed", zap.Error(err))
	}

	zap.L().Info("application stopped successfully")
}

func initGlobalLogger(serviceName string) DeferFunc {
	logger := zap.
		Must(zap.NewProduction()).
		With(zap.String("service", serviceName)) // 필수 필드 추가
	undo := zap.ReplaceGlobals(logger)

	return func() {
		_ = logger.Sync()
		undo()
	}
}

// App app
type App interface {
	Run() error
	Stop() error
}

// Factory app factory
type Factory interface {
	NewApp(config any) (App, error)
	NewConfig() any
}
