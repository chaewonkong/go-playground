package multi_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestMultiErr(t *testing.T) {
	t.Run("TestMultiErr", func(t *testing.T) {
		e1 := errors.New("error 1")
		e2 := errors.New("error 2")

		err := multierr.Combine(e1, e2)

		errors := multierr.Errors(err)
		for idx, e := range errors {
			assert.Equal(t, fmt.Sprintf("error %d", idx+1), e.Error(), "에러 메시지가 일치하지 않습니다")
		}
	})
}
