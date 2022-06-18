package configuration

import (
	"log"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Port                   string  `json:"port"`
	ForwardAdress          string  `json:"forwardAdress"`
	MaxOcupation           int     `json:"maxOcupation"`
	MaxThreshold           int     `json:"threshold"`
	MaxPowerAllowed        float64 `json:"maxPowerAllowed"`
	MaxFrqz                int     `json:"maxFrequenzy"`
	MaxEnergyCostPerPymemo float64 `json:"maxEnergyCostPerPymemo"`
}

func SetConfig(dir string) Configuration {
	filename := dir + "config.json"
	configuration := Configuration{}
	err := gonfig.GetConf(filename, &configuration)
	if err != nil {
		log.Println("No se ha podido leer la configuración en " + dir)
		panic(err)
	}
	return configuration
}
