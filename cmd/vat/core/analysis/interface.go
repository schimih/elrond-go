package analysis

type AnalyzeRunner interface {
	Run()
	GetTargetsList() []Target
}

type Discoverer interface {
	DiscoverNewTargets(existingTargets []Target) (targets []Target)
	IsInterfaceNil() bool
}
