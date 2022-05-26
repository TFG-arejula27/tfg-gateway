package manager

import "log"

const (
	OCUPATION  = 0
	ENERGYCOST = 1
	POWER      = 2
	TIME       = 3
)

type eventMessage struct {
	kind  int16
	value interface{}
}

type State struct {

	//state atributes
	energyCost           float64
	averageExecutionTime float64
	averagePower         float64
	ocupation            int32

	///state communication

	events chan eventMessage
}

func NewState() *State {
	return &State{events: make(chan eventMessage, 100)}
}

func (state *State) Run() {

	func() {
		for {
			m := <-state.events

			//handle events
			switch m.kind {
			case POWER:
				log.Printf("current  power %f", m.value)
				break
			case TIME:
				log.Printf("current  time %f", m.value)
				break

			}
		}
	}()
}

func (state *State) ChangeOcupation(value int) {
	state.events <- eventMessage{OCUPATION, value}
}

func (state *State) ChangeEnergyCost(value float64) {
	state.events <- eventMessage{ENERGYCOST, value}
}

func (state *State) ChangeAveragePower(value float64) {
	state.events <- eventMessage{POWER, value}
}

func (state *State) ChangeExecutionTime(value float64) {
	state.events <- eventMessage{TIME, value}
}
