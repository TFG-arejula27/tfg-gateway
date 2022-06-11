package receptor

import (
	"log"
	"strconv"

	"github.com/arejula27/energy-cluster-manager/internal/manager"
)

type Receptor struct {
	Manager *manager.Manager
}

func NewReceptor(state *manager.Manager) *Receptor {
	return &Receptor{Manager: state}
}

func (r *Receptor) GetCurrentPower() {

	power, err := runPowerstat()
	if err != nil {
		log.Println(err)
		return
	}
	averagePower, err := strconv.ParseFloat(power, 64)
	if err != nil {
		log.Println(err)
		return
	}
	r.Manager.ChangeAveragePower(averagePower)

}
