package analysis

type AnalyzerStub struct {
	AnalyzeCalled func(existingTargets []Target) (targets []Target)
}

func NewAnalyzerStub() (*AnalyzerStub, error) {
	return &AnalyzerStub{}, nil
}
