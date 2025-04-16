package cockroach_test

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorsNew(t *testing.T) {
	err := firstFunction()

	// 1. 에러가 발생했는지 확인
	assert.Error(t, err, "에러가 발생해야 합니다")

	// 2. 에러 메시지 확인
	assert.Equal(t, "err from: 스택 트레이스 테스트용 에러입니다", err.Error())
	assert.Contains(t, err.Error(), "err from:")

	fullErr := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullErr, "secondFunction")
	assert.Contains(t, fullErr, "firstFunction")
	assert.Contains(t, fullErr, "TestErrorsNew")
}

func firstFunction() error {
	return secondFunction()
}

func secondFunction() error {
	// 스택 트레이스를 포함하여 에러 생성
	msg := "스택 트레이스 테스트용 에러입니다"
	return errors.Errorf("err from: %s", msg)
	// return errors.New("스택 트레이스 테스트용 에러입니다")
}

func TestErrorsJoin(t *testing.T) {
	err1 := errors.New("test")
	err2 := errors.New("test2")
	err3 := errors.New("test3")

	err := errors.Join(err1, err2, err3)

	assert.Equal(t, "test\ntest2\ntest3", err.Error())
}
