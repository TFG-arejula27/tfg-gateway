package main

import (
	"github.com/arejula27/energy-cluster-manager/internal/gateway"
)

func main() {
	r := gateway.NewHandler()
	r.Run()

}
