package receptor

import (
	"log"
	"strconv"

	"github.com/arejula27/energy-cluster-manager/internal/manager"
)

type Receptor struct {
	State *manager.State
}

func NewReceptor(state *manager.State) *Receptor {
	return &Receptor{State: state}
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
	r.State.ChangeAveragePower(averagePower)

}
