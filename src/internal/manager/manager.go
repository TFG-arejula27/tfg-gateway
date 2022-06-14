package manager

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	//factors
	energyCost    float64 //en Kwh
	executionTime time.Duration
	averagePower  float64
}

type state struct {
	//state atributes
	sync.Mutex
	stateProperties
}

type restrictions struct {
	//restrictions
	maxAllowedEnergycost float64 //En €/h
	maxAllowedThreshold  int     //[0,255]
}

type Manager struct {

	//strategy
	strategy strategy

	//restrictions
	restrictions restrictions

	///state
	events chan eventMessage
	state  state

	//outpus
	threshold    int
	maxOcupation int
	frequenzy    int
	throghput    float64

	//Log
	log *log.Logger
}

// La monitorización quiźa sea la parte más delicada, necesitamos saber:
// 1) carga de trabajo (como la generamos nosotros, la podemos obtener fácilmente)
// 2) coste de la energía (lo generamos nosotros)
// 3) tiempo de ejecución de las tareas, ahí tú sabes más que yo... ¿Cómo se puede obtener el tiempo de ejecución fácilmente de openfaas / kubernetes?
// 4) consumo energético (esto ya lo sabes hacer)

func NewManager(str strategy, last bool, ocupation int, maxCost float64, maxThreshold int, dir string) *Manager {
	logFile, err := os.OpenFile(dir+"manager.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln("no se ha podido crear archivo de log")
	}

	if err != nil {
		log.Fatal(err)
	}
	logMng := log.New(logFile, "", log.LstdFlags)

	mng := &Manager{
		strategy:     str,
		events:       make(chan eventMessage, 100),
		threshold:    0,
		state:        state{},
		restrictions: restrictions{maxAllowedEnergycost: maxCost, maxAllowedThreshold: maxThreshold},
		maxOcupation: ocupation,
		log:          logMng,
	}

	mng.getFrecuenzy()
	mng.logHeader()
	return mng
}

func (mng *Manager) Run() {

	go mng.simulateEnergyPrice()
	//bucle mape
	time.Sleep(time.Second * 5)
	go func() {
		for {
			time.Sleep(time.Second * 60)
			//monitoring
			mng.state.Lock()
			currentState := stateProperties{
				energyCost:    mng.state.energyCost,
				executionTime: mng.state.executionTime,
				averagePower:  mng.state.averagePower,
			}
			mng.state.Unlock()

			//si gastamos más
			if currentState.averagePower*currentState.energyCost > mng.restrictions.maxAllowedEnergycost {
				//analyze and planing
				decision := mng.strategy.takeDecision(currentState, mng.restrictions)
				//execute
				mng.doDecision(decision)
				log.Println("Decision tomada", decision.frecuenzy, decision.ocupation, decision.thrshold)

			}
			//si estamos en la potencia correcta no cambiar nada
			mng.logCurrentStatus()

		}
	}()

}

func (mng *Manager) simulateEnergyPrice() {
	e := 0.9
	for {
		mng.state.energyCost = e
		time.Sleep(61 * time.Second)
		e += 0.4

	}

}

func (mng *Manager) logHeader() {
	var line string
	//threshold
	line += "threshold "
	//ocupación
	line += "ocupation "
	//frecuencia
	line += "frequenzy "
	//potencia
	line += "power "
	//coste energético
	line += "energyCost "
	//throghtput
	line += "throghtput "

	mng.log.Println(line)

}
func (mng *Manager) logCurrentStatus() {

	var line string
	//threshold
	line += strconv.Itoa(mng.threshold) + " "
	//ocupación
	line += strconv.Itoa(mng.maxOcupation) + " "
	//frecuencia
	line += strconv.Itoa(mng.frequenzy) + " "
	//potencia
	line += strconv.FormatFloat(mng.state.averagePower, 'f', 4, 64) + " "
	//coste energético
	line += strconv.FormatFloat(mng.state.energyCost, 'f', 4, 64) + " "
	//throghtput

	line += strconv.FormatFloat(mng.throghput, 'f', 4, 64) + " "

	mng.log.Println(line)
}

func (mng *Manager) doDecision(d decision) error {
	mng.state.Lock()
	defer mng.state.Unlock()
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
	mng.state.Lock()
	mng.state.averagePower = value
	mng.state.Unlock()
	log.Println("current  power ", value)
}

func (mng *Manager) ChangeExecutionTime(value time.Duration) {
	mng.state.Lock()
	mng.state.executionTime = value
	mng.state.Unlock()
	log.Println("current  time ", value.Seconds())
}

func (mng *Manager) ChangeEnergyCost(value float64) {
	mng.state.Lock()
	mng.state.energyCost = value
	mng.state.Unlock()
	log.Println("current  energy cost ", value)
}

//Resource management functions
func (mng *Manager) setFrecuenzy(frequenzy int) error {
	mng.frequenzy = frequenzy
	freq := strconv.Itoa(frequenzy)
	cmd := exec.Command("cpupower", "frequency-set", "--freq", freq)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func (mng *Manager) getFrecuenzy() error {
	cmd := exec.Command("cpupower", "frequency-info", "-f")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(output), "\n")
	line := strings.Join(strings.Fields(lines[1]), " ")
	freq := strings.Split(line, " ")[3]

	f, _ := strconv.Atoi(freq)
	mng.frequenzy = f
	return nil
}

func (mng *Manager) setMaxOcupation(ocupation int) {

	mng.maxOcupation = ocupation

}

func (mng *Manager) setThreshold(threshold int) {
	mng.threshold = threshold

}

func (mng *Manager) ChangeThroughput(throghput float64) {
	mng.throghput = throghput

}

//Getter functions outputs

func (mng *Manager) GetMaxOcupation() int {
	mng.state.Lock()
	defer mng.state.Unlock()
	return mng.maxOcupation
}

func (mng *Manager) GetThreshold() int {
	mng.state.Lock()
	defer mng.state.Unlock()
	return mng.threshold
}
