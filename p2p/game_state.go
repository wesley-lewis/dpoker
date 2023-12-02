package p2p 

type Round uint32

const (
	Dealing Round = iota
	Preflop
	Flop
	Turn
	River
)

type GameState struct {
	isDealer			bool 
	Round					Round
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) loop() {
	for {
		select {

		}
	}
}
