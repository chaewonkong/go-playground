package player

const defaultScore = 1500.0

// Player 경기 참가자 인터페이스
type Player interface {
	// Score 점수 getter
	Score() float64

	// SetScore 점수 setter
	SetScore(score float64)
}

type player struct {
	score float64
	// consecutiveWins int
}

func New() Player {
	return &player{
		score: defaultScore,
	}
}

func (p *player) Score() float64 {
	return p.score
}

func (p *player) SetScore(score float64) {
	p.score = score
}
