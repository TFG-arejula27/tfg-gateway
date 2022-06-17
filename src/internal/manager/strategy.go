package manager

const (
	FORWARD = 0
	LOCAL   = 1

	MODEL_DATA_DIR = "model/"
)

type decision struct {
	frecuenzy int
	thrshold  int
	ocupation int
}

type stats struct {
	energyWasted float64
	time         float64
}
type strategy interface {
	takeDecision(state state, restrictions restrictions) (decision, stats)
}
