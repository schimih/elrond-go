package analysis

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/stretchr/testify/assert"
)

type FakeDiscoverer struct {
}

func (fd *FakeDiscoverer) DiscoverNewTargets(existingTargets []Target) (targets []Target) {
	targets = existingTargets

	return
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *FakeDiscoverer) IsInterfaceNil() bool {
	return d == nil
}
func TestNewAnalyzer(t *testing.T) {
	analysisType := "tralala"
	fd := &FakeDiscoverer{}
	na, err := NewAnalyzer(fd, analysisType)
	assert.False(t, check.IfNil(na))
	assert.Nil(t, err)
}
