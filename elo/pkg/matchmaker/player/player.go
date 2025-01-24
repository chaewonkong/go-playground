package player

import (
	"time"
)

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

// Ticket ticket
type Ticket struct {
	// Matchmaking rating
	MMR float64

	// Ticket created at
	CreatedAt time.Time
}

type Pool struct {
	tickets []Ticket
}

func (p *Pool) Add(t Ticket) {

}

// sort()
// MMR + waitingTime의 계산결과를 sorting score로 pool을 정렬
// Tick sort가 진행되는 시간 단위 (1s)
// 매 Tick마다 sort를 진행하고 match를 찾음?
