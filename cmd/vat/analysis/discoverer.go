package analysis

import (
	"strings"

	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/elrond-go/cmd/vat/core"
)

type P2pDiscoverer struct {
	messenger p2p.Messenger
}

func NewP2pDiscoverer(messenger p2p.Messenger) *P2pDiscoverer {
	return &P2pDiscoverer{
		messenger: messenger,
	}
}

func (d *P2pDiscoverer) DiscoverNewTargets(targetsDiscoveredLastRound []DiscoveredTarget) (discoveredTargets []DiscoveredTarget) {
	discoveredTargets = targetsDiscoveredLastRound
	currentlyConnectedTargets := d.messenger.ConnectedAddresses()

	for idx, address := range currentlyConnectedTargets {
		targetAddress := strings.Split(address, "/")
		target := MakeTarget(uint(idx), targetAddress[1], targetAddress[2], targetAddress[4], core.NEW)
		if !containsTarget(discoveredTargets, target) {
			discoveredTargets = append(discoveredTargets, target)
		}
	}
	return
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *P2pDiscoverer) IsInterfaceNil() bool {
	return d == nil
}
