package analysis

type DiscovererStub struct {
	DiscoverNewTargetsCalled func(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget)
}

func NewDiscovererStub() *DiscovererStub {
	return &DiscovererStub{}
}

func (stub *DiscovererStub) DiscoverNewTargets(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget) {
	if stub.DiscoverNewTargetsCalled != nil {
		return stub.DiscoverNewTargetsCalled(existingTargets)
	}

	return make([]DiscoveredTarget, 0)
}

// IsInterfaceNil returns true if there is no value under the interface
func (stub *DiscovererStub) IsInterfaceNil() bool {
	return stub == nil
}
