package main

import (
	"time"

	"github.com/arejula27/energy-cluster-manager/internal/gateway"
	"github.com/arejula27/energy-cluster-manager/internal/manager"
	"github.com/arejula27/energy-cluster-manager/internal/receptor"
)

func main() {
	state := manager.NewState()
	r := receptor.NewReceptor(state)
	go runManager(r)

	go state.Run()
	runHandler(r)

}

func runManager(receptor *receptor.Receptor) {
	for {
		//cada dos minutos tomar metrica
		time.Sleep(time.Second * 120)
		receptor.GetCurrentPower()
	}

}

func runHandler(receptor *receptor.Receptor) {
	r := gateway.NewHandler(receptor)
	r.Gin.Run()

}
