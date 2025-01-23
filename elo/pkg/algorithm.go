package pkg

import (
	"elo/pkg/algorithm"
)

const (
	divisor  = 400
	defaultK = 60
)

// Algorithm 알고리즘
type Algorithm struct {
	K int
}

// MatchingPool 매칭 풀
type MatchingPool struct{}

// New 생성자
func New() Algorithm {
	return Algorithm{
		K: defaultK,
	}
}

func calculate(scoreA, scoreB, result float64) (resultA, resultB float64) {
	expectedA := algorithm.GetWinProb(scoreA, scoreB)
	expectedB := 1 - expectedA

	resultA = scoreA + (defaultK * (result - expectedA))
	resultB = scoreB + (defaultK * (result - expectedB))

	return
}

/*
1. MMR
2. MatchMaking
*/
