package manager

import (
	"log"
	"os/exec"
	"strconv"
)

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

	//outpus
	Threshold int
}

func NewState() *State {
	return &State{events: make(chan eventMessage, 100), Threshold: 125}
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
			case OCUPATION:
				log.Printf("current  ocupation %f", m.value)
			case ENERGYCOST:
				log.Printf("current  energy cost %f", m.value)

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

func setFrecuenzy(frequenzy int) error {
	freq := strconv.Itoa(frequenzy)
	cmd := exec.Command("cpupower", "frequency-set", "--freq", freq)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}
