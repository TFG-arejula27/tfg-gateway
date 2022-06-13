package manager

const (
	FORWARD = 0
	LOCAL   = 1
)

type decision struct {
	frecuenzy int
	thrshold  int
	ocupation int
}
type strategy interface {
	takeDecision(s stateProperties) decision
}
