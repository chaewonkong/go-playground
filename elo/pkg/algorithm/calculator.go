package algorithm

const defaultK = 20

// CalculateScore 승률과 현재 점수를 바탕으로 새 점수 계산
func CalculateScore(currScore, winProb float64, result float64) float64 {
	delta := defaultK * (result - winProb)

	return currScore + delta
}
