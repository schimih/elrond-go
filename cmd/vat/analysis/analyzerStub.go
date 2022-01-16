package analysis

type AnalyzerStub struct {
	AnalyzeCalled func(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget)
}

func NewAnalyzerStub() (*AnalyzerStub, error) {
	return &AnalyzerStub{}, nil
}
