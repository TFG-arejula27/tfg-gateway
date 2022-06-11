package manager

const (
	FORWARD = 0
	LOCAL   = 1
)

type decision struct {
	value int
}
type strategy interface {
	takeDecision(s state) decision
}

type dumb struct {
	maxOcupation int
}

func NewDumbStrategy(max int) *dumb {
	return &dumb{maxOcupation: max}
}

func (str *dumb) takeDecision(s state) decision {
	return decision{value: LOCAL}
	if s.last {
		return decision{value: LOCAL}
	} else if s.ocupation == int(str.maxOcupation) {
		return decision{value: FORWARD}
	}

	return decision{value: LOCAL}

}
