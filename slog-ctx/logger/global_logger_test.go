package logger_test

import (
	"sync"
	"testing"

	. "slog-ctx/logger"

	"github.com/stretchr/testify/assert"
)

func TestGlobalLogger(t *testing.T) {
	t.Run("Init하지 않은 경우에도 logger는 nil이 아님", func(t *testing.T) {
		assert.NotNil(t, G())
	})

	t.Run("동시성 테스트 with go test -race", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(2)

			go func() {
				defer wg.Done()
				InitGlobalLogger()
			}()
			go func() {
				defer wg.Done()
				_ = G()
			}()
		}
		wg.Wait()
	})
}
