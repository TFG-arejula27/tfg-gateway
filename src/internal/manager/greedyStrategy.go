package manager

type GreedyStratey struct {
	maxOcupation int
}

func NewGreedyStratey(ocupation int) *GreedyStratey {
	return &GreedyStratey{maxOcupation: ocupation}
}

func (strat *GreedyStratey) takeDecision(s stateProperties) decision {
	return decision{}
}
