package manager

import (
	"fmt"
	"log"
	"os"
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

type state struct {
	//state atributes
	executionTime time.Duration
	averagePower  float64
	energyPrice   float64 //Kwh
}

type restrictions struct {
	//restrictions
	maxAllowedPower         float64 //En Watts
	maxAllowedThreshold     int     //[0,255]
	maxAllowedCostPerPymemo float64 //en €
}

type Manager struct {
	sync.Mutex

	//strategy
	strategy strategy

	//restrictions
	restrictions restrictions

	///state
	state state
	stats stats
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

func NewManager(str strategy, last bool, ocupation int, maxAllowedPower float64, maxThreshold int, dir string, maxFrqz int, maxAllowedCostPerPymemo float64) *Manager {
	logFile, err := os.OpenFile(dir+"manager.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln("no se ha podido crear archivo de log")
	}

	if err != nil {
		log.Fatal(err)
	}
	logMng := log.New(logFile, "", log.LstdFlags)

	mng := &Manager{
		strategy:  str,
		threshold: 0,
		state:     state{},
		restrictions: restrictions{
			maxAllowedPower:         maxAllowedPower,
			maxAllowedThreshold:     maxThreshold,
			maxAllowedCostPerPymemo: maxAllowedCostPerPymemo},
		maxOcupation: ocupation,
		log:          logMng,
	}

	mng.setFrecuenzy(maxFrqz)
	mng.logHeader()
	return mng
}

func (mng *Manager) Run() {

	go mng.simulateEnergyPrice()
	//bucle mape
	time.Sleep(time.Second * 5)
	go func() {
		for {
			fmt.Println("60 s until text mape iteration")
			time.Sleep(time.Second * 60)
			//monitoring
			mng.Lock()
			currentState := state{
				executionTime: mng.state.executionTime,
				averagePower:  mng.state.averagePower,
			}
			mng.Unlock()

			//si gastamos más
			if currentState.averagePower > mng.restrictions.maxAllowedPower {
				//analyze and planing
				decision, stats := mng.strategy.takeDecision(currentState, mng.restrictions)
				mng.stats = stats
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
	//todo ponerlo bien
	add := float64(-3)
	for {
		if mng.restrictions.maxAllowedPower <= 45 {
			add = float64(3)
		}
		mng.restrictions.maxAllowedPower += add
		time.Sleep(241 * time.Second)

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
	line += "power"
	//coste energético
	line += "energyCost "
	//throghtput
	line += "throghtput"
	//evolución energía un pymemo
	line += "pymemoEnergy"
	//lamda

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
	//coste más allowed
	line += strconv.FormatFloat(mng.restrictions.maxAllowedPower, 'f', 4, 64) + " "
	//throghtput
	line += strconv.FormatFloat(mng.throghput, 'f', 4, 64) + " "
	//evolución energía un pymemo
	line += strconv.FormatFloat(mng.state.executionTime.Seconds()*mng.state.averagePower, 'f', 4, 64) + " "
	mng.log.Println(line)
}

func (mng *Manager) doDecision(d decision) error {
	mng.Lock()
	defer mng.Unlock()
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
	mng.Lock()
	mng.state.averagePower = value
	mng.Unlock()
	log.Println("current  power ", value)
}

func (mng *Manager) ChangeExecutionTime(value time.Duration) {
	mng.Lock()
	mng.state.executionTime = value
	mng.Unlock()
	log.Println("current  time ", value.Seconds())
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
	return nil
}

func (mng *Manager) setMaxOcupation(ocupation int) {
	mng.Lock()
	defer mng.Unlock()
	mng.maxOcupation = ocupation

}

func (mng *Manager) setThreshold(threshold int) {
	mng.Lock()
	defer mng.Unlock()
	mng.threshold = threshold

}

func (mng *Manager) ChangeThroughput(throghput float64) {
	mng.Lock()
	defer mng.Unlock()
	mng.throghput = throghput

}

//Getter functions outputs

func (mng *Manager) GetMaxOcupation() int {
	mng.Lock()
	defer mng.Unlock()
	return mng.maxOcupation
}

func (mng *Manager) GetThreshold() int {
	mng.Lock()
	defer mng.Unlock()
	return mng.threshold
}
