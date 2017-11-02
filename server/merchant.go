package server

const (
	BUY = 0
	BUYING = 1
	SELL = 2
	SELLING = 3
)

type merchant struct {
	ch chan string
	state int
}

func (m *merchant) Listen() {
	for {
		switch m.state {
		case BUY:

		}
	}
}
