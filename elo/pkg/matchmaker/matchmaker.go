package matchmaker

import (
	"errors"

	"github.com/chaewonkong/go-playground/elo/pkg/matchmaker/player"
)

type Matchmaker interface {
}

type Ticket struct {
	ID     string
	Player player.Player
}

type Match struct {
	Tickets []Ticket
}

type LocalMatchmaker struct {
	pool Pool
}

type Pool struct {
	tickets []Ticket
}

func (p *Pool) Add(tgt Ticket) {
	tickets := make([]Ticket, 0, len(p.tickets)+1)
	for _, t := range p.tickets {
		s := t.Player.Score()
		if s > tgt.Player.Score() {
			tickets = append(tickets, tgt)
		}
		tickets = append(tickets, t)
	}

	p.tickets = tickets
}

func (p *Pool) Poll() (t Ticket, err error) {
	if len(p.tickets) < 1 {
		err = errors.New("No item left to poll")
		return
	}
	t = p.tickets[0]
	p.tickets = p.tickets[1:]
	return
}

func (p *Pool) Length() int {
	return len(p.tickets)
}

func (m *LocalMatchmaker) FindMatch(matchSize int) (match Match, ok bool) {
	for i := 0; i < matchSize; i++ {
		tkt, err := m.pool.Poll()
		if err != nil {
			return
		}
		match.Tickets = append(match.Tickets, tkt)
	}

	ok = true
	return
}
