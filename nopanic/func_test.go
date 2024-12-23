package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilNotPanic(t *testing.T) {
	t.Run("interface가 nil로 주입되어도 panic하지 않음", func(t *testing.T) {
		c := NewChildStruct()
		assert.Nil(t, c)
		p := NewParentStruct(c)
		assert.NotNil(t, p)
		assert.NotPanics(t, func() {
			p.Hello()
		})
	})
}
