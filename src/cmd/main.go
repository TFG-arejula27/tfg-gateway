package main

import (
	"github.com/arejula27/energy-cluster-manager/internal/gateway"
	"github.com/arejula27/energy-cluster-manager/internal/manager"
	"github.com/arejula27/energy-cluster-manager/internal/receptor"
)

func main() {

	runManager()
	//runHandler()

}

func runManager() {
	state := manager.NewState()
	state.Run()

	r := receptor.NewReceptor(state)
	r.GetCurrentPower()
}

func runHandler() {
	r := gateway.NewHandler()
	r.Run()

}
