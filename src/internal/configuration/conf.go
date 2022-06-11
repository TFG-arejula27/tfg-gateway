package configuration

import "github.com/tkanos/gonfig"

type Configuration struct {
	Port          string `json:"port"`
	ForwardAdress string `json:"forwardAdress"`
	MaxOcupation  int    `json:"maxOcupation"`
}

func SetConfig() Configuration {
	filename := "./config.json"
	configuration := Configuration{}
	gonfig.GetConf(filename, &configuration)
	return configuration
}
