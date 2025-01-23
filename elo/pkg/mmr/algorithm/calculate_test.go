package algorithm_test

import (
	"elo/pkg/mmr/algorithm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetScoreFromGameResult(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// given
		pvtScore := 2000.0
		gameResult := 1.0
		winProb := 0.64

		expected := 2021.6
		// when
		actual := algorithm.CalculateScore(pvtScore, winProb, gameResult)

		// then
		assert.Equal(t, expected, actual)
	})
}
