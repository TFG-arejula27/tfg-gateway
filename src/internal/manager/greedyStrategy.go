package manager

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/gocarina/gocsv"
)

type ModelElement struct {
	decision   decision
	power      float64
	energy     float64
	throghtput float64
}

type GreedyStratey struct {
	maxOcupation int
	model        []ModelElement
}

func NewGreedyStratey(ocupation int, modelDir string) *GreedyStratey {

	strat := GreedyStratey{maxOcupation: ocupation}

	err := strat.loadData(modelDir)
	if err != nil {
		log.Println("Can't load model from files")
		panic(err)
	}
	return &strat
}

func (strat *GreedyStratey) loadData(dir string) error {

	type Log struct {
		Message   string  `csv:"Message"`
		Power     float64 `csv:"Average power(Watts)"`
		Time      float64 `csv:"time"`
		Threshold int     `csv:"threshold"`
		Frecuenzy int     `csv:"Frecuenzy"`
		Energy    float64 `csv:"Energy"`
	}

	for i := 1; i <= strat.maxOcupation; i++ {
		fileName := dir + MODEL_DATA_DIR + strconv.Itoa(i) + ".csv"
		file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Println("Can't load model for ocupation = ", i)
			continue

		}

		gocsv.SetCSVReader(func(out io.Reader) gocsv.CSVReader {
			reader := csv.NewReader(out)
			reader.Comma = ' '
			return reader

		})

		logs := []*Log{}

		if err := gocsv.UnmarshalFile(file, &logs); err != nil { // Load clients from file
			panic(err)
		}
		//leer el fichero
		for _, log := range logs {
			throghtput := float64(i) / log.Time

			newElement := ModelElement{
				power:      log.Power,
				throghtput: throghtput,
				energy:     log.Energy,
				decision: decision{
					frecuenzy: log.Frecuenzy,
					thrshold:  log.Threshold,
					ocupation: i,
				}}
			strat.model = append(strat.model, newElement)

		}

	}
	//cargados todos los datos los ordenamos por throughput
	sort.Slice(strat.model, func(i, j int) bool {

		return strat.model[i].throghtput > strat.model[j].throghtput

	})
	for _, v := range strat.model {
		if v.decision.frecuenzy == 1600000 && v.decision.ocupation == 5 {
			fmt.Println(v)
		}

	}

	return nil

}

func (strat *GreedyStratey) takeDecision(state state, restrictions restrictions) (decision, stats) {
	conjunto := strat.model

	//select first
	for _, modelElement := range conjunto {
		if checkRestrctions(state, modelElement, restrictions) {
			//es solución
			stats := stats{energyWasted: modelElement.energy}
			return modelElement.decision, stats
		}

	}

	return decision{ocupation: 0}, stats{}

}

func checkRestrctions(s state, e ModelElement, r restrictions) bool {

	if beatsMaxCostPerPymemo(s, e, r) {
		return false
	}
	//si el threshold es mayor del permitido
	if e.decision.thrshold > r.maxAllowedThreshold {
		return false
	}
	return true

}

func beatsMaxPower(s state, e ModelElement, r restrictions) bool {
	//si el coste energético es mayor del permitido
	log.Println("modelo", e.power, "allowed", r.maxAllowedPower)
	return e.power > r.maxAllowedPower

}

func beatsMaxCostPerPymemo(s state, e ModelElement, r restrictions) bool {
	pymemoEnergyCost := e.energy / float64(e.decision.ocupation) //J
	maxCost := pymemoEnergyCost / 3600000 * s.energyPrice
	if maxCost > r.maxAllowedCostPerPymemo {
		return false
	}
	return true

}

func beatsMeanMaxCostPerPymemo(s state, e ModelElement, r restrictions) bool {
	return true

}
