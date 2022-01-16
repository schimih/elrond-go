package analysis

type DiscovererStub struct {
	DiscoverNewTargetsCalled func(existingTargets []Target) (targets []Target)
}

func NewDiscovererStub() *DiscovererStub {
	return &DiscovererStub{}
}

func (stub *DiscovererStub) DiscoverNewTargets(existingTargets []Target) (targets []Target) {
	if stub.DiscoverNewTargetsCalled != nil {
		return stub.DiscoverNewTargetsCalled(existingTargets)
	}

	return make([]Target, 0)
}

// IsInterfaceNil returns true if there is no value under the interface
func (stub *DiscovererStub) IsInterfaceNil() bool {
	return stub == nil
}
