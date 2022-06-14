package main

import (
	"flag"
	"fmt"

	"github.com/arejula27/energy-cluster-manager/internal/configuration"
	"github.com/arejula27/energy-cluster-manager/internal/gateway"
	"github.com/arejula27/energy-cluster-manager/internal/manager"
	"github.com/arejula27/energy-cluster-manager/internal/receptor"
)

func main() {

	conf, dirConf := initConf()

	strategy := manager.NewGreedyStratey(conf.MaxOcupation, dirConf)

	last := conf.ForwardAdress == ""
	manager := manager.NewManager(strategy, last, conf.MaxOcupation, conf.MaxEnergyCost, conf.MaxThreshold, dirConf, conf.MaxFrqz)
	fmt.Println(conf)
	r := receptor.NewReceptor(manager)
	go runReceptor(r)

	go manager.Run()
	runHandler(manager, conf)

}

func initConf() (configuration.Configuration, string) {

	dir := flag.String("file", "~/.rscManager/", "set the directory with the configuration")
	flag.Parse()
	return configuration.SetConfig(*dir), *dir
}

func runReceptor(receptor *receptor.Receptor) {
	for {
		//cada dos minutos tomar metrica

		receptor.GetCurrentPower()
		//time.Sleep(time.Second * 12)
	}

}

func runHandler(manager *manager.Manager, conf configuration.Configuration) {
	r := gateway.NewHandler(manager, conf)
	r.Gin.Run(":" + conf.Port)

}
