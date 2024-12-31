package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const errFilterNotImplemented = "filter not implemented"

// FilterNotImplementedError 필터가 구현되지 않은 경우 발생하는 에러
type FilterNotImplementedError struct {
	msg string
}

// NewFilterNotImplementedError 생성자
func NewFilterNotImplementedError() error {
	return &FilterNotImplementedError{errFilterNotImplemented}
}

func (e *FilterNotImplementedError) Error() string {
	return e.msg
}

// FilterOptions 필터 생성을 위한 옵션값
type FilterOptions struct {
	NumTeams    int
	MaxTeamSize int
}

func TestErrorIs(t *testing.T) {
	// given
	// when
	var err error = NewFilterNotImplementedError()
	// then
	_, ok := err.(*FilterNotImplementedError)
	assert.True(t, ok)
}
