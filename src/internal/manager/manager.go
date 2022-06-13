package manager

import (
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
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

type stateProperties struct {
	energyCost           float64
	averageExecutionTime time.Duration
	averagePower         float64
}

type state struct {
	//state atributes
	sync.Mutex
	stateProperties
}

type Manager struct {

	//strategy
	strategy strategy

	///state communication

	events chan eventMessage

	state state

	//outpus
	threshold    int
	maxOcupation int
}

// La monitorización quiźa sea la parte más delicada, necesitamos saber:
// 1) carga de trabajo (como la generamos nosotros, la podemos obtener fácilmente)
// 2) coste de la energía (lo generamos nosotros)
// 3) tiempo de ejecución de las tareas, ahí tú sabes más que yo... ¿Cómo se puede obtener el tiempo de ejecución fácilmente de openfaas / kubernetes?
// 4) consumo energético (esto ya lo sabes hacer)

func NewManager(str strategy, last bool, ocupation int) *Manager {
	return &Manager{
		strategy:     str,
		events:       make(chan eventMessage, 100),
		threshold:    0,
		state:        state{},
		maxOcupation: ocupation,
	}
}

func (mng *Manager) Run() {
	//bucle mape
	go func() {
		for {
			//monitoring
			mng.state.Lock()
			currentState := stateProperties{
				energyCost:           mng.state.energyCost,
				averageExecutionTime: mng.state.averageExecutionTime,
				averagePower:         mng.state.averagePower,
			}
			mng.state.Unlock()
			//analyze and planing
			decision := mng.strategy.takeDecision(currentState)
			//execute
			mng.doDecision(decision)
			time.Sleep(time.Second * 10)
		}
	}()

	func() {
		for {
			m := <-mng.events
			//handle events
			switch m.kind {
			case POWER:
				mng.state.Lock()
				mng.state.averagePower = m.value.(float64)
				mng.state.Unlock()
				log.Printf("current  power %f", m.value)
				break
			case TIME:
				mng.state.Lock()
				mng.state.averageExecutionTime = m.value.(time.Duration)
				mng.state.Unlock()
				log.Printf("current  time %f", m.value)
				break

			case ENERGYCOST:
				mng.state.Lock()
				mng.state.energyCost = m.value.(float64)
				mng.state.Unlock()
				log.Printf("current  energy cost %f", m.value)

			}
		}
	}()
}

func (mng *Manager) doDecision(d decision) error {

	return nil //TODO quitar
	err := mng.setFrecuenzy(d.frecuenzy)
	if err != nil {
		return err
	}
	mng.setMaxOcupation(d.ocupation)
	mng.setThreshold(d.thrshold)

	return nil

}

//monitoring functions
func (mng *Manager) ChangeAveragePower(value float64) {
	mng.events <- eventMessage{POWER, value}
}

func (mng *Manager) ChangeExecutionTime(value time.Duration) {
	mng.events <- eventMessage{TIME, value}
}

func (mng *Manager) ChangeEnergyCost(value float64) {
	mng.events <- eventMessage{ENERGYCOST, value}
}

//Resource management functions
func (mng *Manager) setFrecuenzy(frequenzy int) error {
	freq := strconv.Itoa(frequenzy)
	cmd := exec.Command("cpupower", "frequency-set", "--freq", freq)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func (mng *Manager) setMaxOcupation(ocupation int) {
	mng.maxOcupation = ocupation

}

func (mng *Manager) setThreshold(threshold int) {
	mng.threshold = threshold

}

//Getter functions outputs

func (mng *Manager) GetMaxOcupation() int {
	return mng.maxOcupation
}

func (mng *Manager) GetThreshold() int {
	return mng.threshold
}
