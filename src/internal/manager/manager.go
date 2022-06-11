package manager

import (
	"log"
	"os/exec"
	"strconv"
)

const (
	ENERGYCOST = 1
	POWER      = 2
	TIME       = 3
	EVAL       = 4
)

type eventMessage struct {
	kind  int16
	value interface{}
}

type state struct {
	//state atributes
	energyCost           []float64
	averageExecutionTime []float64
	averagePower         []float64
	ocupation            int
	last                 bool
}

type Manager struct {

	//strategy
	strategy strategy

	///state communication

	events      chan eventMessage
	replies     chan *chan (bool)
	ocupationCh chan int

	state state

	//outpus
	Threshold    int
	Forward      bool
	maxOcupation int
}

// La monitorización quiźa sea la parte más delicada, necesitamos saber:
// 1) carga de trabajo (como la generamos nosotros, la podemos obtener fácilmente)
// 2) coste de la energía (lo generamos nosotros)
// 3) tiempo de ejecución de las tareas, ahí tú sabes más que yo... ¿Cómo se puede obtener el tiempo de ejecución fácilmente de openfaas / kubernetes?
// 4) consumo energético (esto ya lo sabes hacer)

func NewManager(str strategy, last bool, ocupation int) *Manager {
	return &Manager{
		strategy:    str,
		events:      make(chan eventMessage, 100),
		ocupationCh: make(chan int),
		Threshold:   0, Forward: false,
		replies:      make(chan *chan bool, 100),
		maxOcupation: ocupation,
	}
}

func (mng *Manager) Run() {

	go func() {
		for {
			o := <-mng.ocupationCh
			mng.state.ocupation += o
			log.Printf("current  ocupation %d", mng.state.ocupation)

			if mng.maxOcupation >= mng.state.ocupation && o < 0 {
				//si es menor que cero liberar un encolado
				free := <-mng.replies
				*free <- true

			}

		}
	}()

	func() {
		for {
			m := <-mng.events

			//handle events
			switch m.kind {
			case POWER:
				log.Printf("current  power %f", m.value)
				break
			case TIME:
				log.Printf("current  time %f", m.value)
				break

			case ENERGYCOST:
				log.Printf("current  energy cost %f", m.value)

			}
		}
	}()
}

func (mng *Manager) doDecision(d decision) {

	switch d.value {
	case FORWARD:
		mng.Forward = true
	case LOCAL:
		mng.Forward = false

	}

}
func (msg *Manager) Eval() {
	decision := msg.strategy.takeDecision(msg.state)
	msg.doDecision(decision)
}

func (mng *Manager) AddExecution(ch *chan (bool)) {
	if mng.maxOcupation > mng.state.ocupation {
		*ch <- true
	} else {
		mng.replies <- ch
	}

}

func (state *Manager) ChangeOcupation(value int) {
	state.ocupationCh <- value
}

func (state *Manager) ChangeEnergyCost(value float64) {
	state.events <- eventMessage{ENERGYCOST, value}
}

func (state *Manager) ChangeAveragePower(value float64) {
	state.events <- eventMessage{POWER, value}
}

func (state *Manager) ChangeExecutionTime(value float64) {
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
