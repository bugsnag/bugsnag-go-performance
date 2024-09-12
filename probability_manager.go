package bugsnagperformance

type probabilityManager struct {
	probability float64
}

func CreateProbabilityManager() *probabilityManager {
	// TODO - implement probability manager
	return &probabilityManager{
		probability: 0.5,
	}
}

func (pm *probabilityManager) getProbability() float64 {
	return pm.probability
}

func (pm *probabilityManager) setProbability(probability float64) {
	pm.probability = probability
}