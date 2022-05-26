package receptor

import (
	"log"
	"strconv"

	"github.com/arejula27/energy-cluster-manager/internal/manager"
)

type receptor struct {
	state *manager.State
}

func NewReceptor(state *manager.State) *receptor {
	return &receptor{}
}

func (r *receptor) GetCurrentPower() {

	power, err := runPowerstat()
	if err != nil {
		log.Println(err)
		return
	}
	averagePower, err := strconv.ParseFloat(power.Averge.Power, 64)
	if err != nil {
		log.Println(err)
		return
	}
	r.state.ChangeAveragePower(averagePower)

}
