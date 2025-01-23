package algorithm_test

import (
	"elo/pkg/algorithm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWinProb(t *testing.T) {
	t.Run("승률을 정상적으로 계산", func(t *testing.T) {
		// given
		pivotScore := 1500.0
		targetScore := 1400.0
		expected := 0.64

		// when
		actual := algorithm.GetWinProb(pivotScore, targetScore)
		assert.Equal(t, expected, actual)
	})

	t.Run("Elo Score가 동점일 때 승률은 0.5", func(t *testing.T) {
		// given
		pivotScore := 1500.0
		targetScore := 1500.0
		expected := 0.5

		// when
		actual := algorithm.GetWinProb(pivotScore, targetScore)
		assert.Equal(t, expected, actual)
	})
}
