package app_test

import (
	. "runapplication/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err := ReadConfig(nil)
			assert.Error(t, err)
		})
	})
}
