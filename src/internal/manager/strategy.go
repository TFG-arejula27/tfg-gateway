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
type strategy interface {
	takeDecision(state stateProperties, restrictions restrictions) decision
}
