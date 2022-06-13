package main

import (
	"fmt"
	"time"

	"github.com/arejula27/energy-cluster-manager/internal/configuration"
	"github.com/arejula27/energy-cluster-manager/internal/gateway"
	"github.com/arejula27/energy-cluster-manager/internal/manager"
	"github.com/arejula27/energy-cluster-manager/internal/receptor"
)

func main() {
	conf := configuration.SetConfig()

	strategy := manager.NewGreedyStratey(conf.MaxOcupation)

	last := conf.ForwardAdress == ""
	manager := manager.NewManager(strategy, last, conf.MaxOcupation)
	fmt.Println(conf)
	r := receptor.NewReceptor(manager)
	go runReceptor(r)

	go manager.Run()
	runHandler(manager, conf)

}

func runReceptor(receptor *receptor.Receptor) {
	for {
		//cada dos minutos tomar metrica
		time.Sleep(time.Second * 120)
		receptor.GetCurrentPower()
	}

}

func runHandler(manager *manager.Manager, conf configuration.Configuration) {
	r := gateway.NewHandler(manager, conf)
	r.Gin.Run(":" + conf.Port)

}
