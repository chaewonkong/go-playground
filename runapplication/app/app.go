package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 환경변수에 저장된 config 파일의 path
const envConfigFilePath = "CONFIG_PATH"

// DeferFunc 정리 함수
type DeferFunc func()

// NewApplication Application을 생성하는 함수
type NewApplication func(readConfig func(config any) error) (App, error)

// RunApplication Application을 실행하고 Signal을 받아 종료 처리를 담당하는 함수
func RunApplication(serviceName string, newApplication NewApplication) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT) // SIGTERM: k8s가 pod를 stop할 때 발생

	// global logger
	cleanUpLogger := initGlobalLogger(serviceName)
	defer cleanUpLogger()

	// Application 생성
	app, err := newApplication(readConfig)
	if err != nil {
		zap.L().Fatal("application init failed", zap.Error(err))
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

// env var에서 config를 읽어온다
func readConfig(cfg any) error {
	configPath := os.Getenv(envConfigFilePath) // config.yml path는 environ으로 관리

	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// App app
type App interface {
	Run() error
	Stop() error
}
