package bugsnagperformance

type probabilityManager struct {
	probability float64
}

func createProbabilityManager() *probabilityManager {
	// TODO - implement probability manager
	return &probabilityManager{
		probability: 1.0,
	}
}

func (pm *probabilityManager) getProbability() float64 {
	return pm.probability
}

func (pm *probabilityManager) setProbability(probability float64) {
	pm.probability = probability
}
