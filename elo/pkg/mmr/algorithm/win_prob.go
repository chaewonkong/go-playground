package algorithm

import "math"

// GetWinProb pivot의 승률 예측
func GetWinProb(pivotScore, targetScore float64) float64 {
	exponent := (targetScore - pivotScore) / 400
	divisor := 1 + math.Pow(10, exponent)
	result := 1 / divisor

	return math.Floor(result*100) / 100
}
