package analysis

import (
	"strings"

	"github.com/ElrondNetwork/elrond-go/p2p"
)

type P2pDiscoverer struct {
	messenger p2p.Messenger
}

func NewP2pDiscoverer(messenger p2p.Messenger) *P2pDiscoverer {
	return &P2pDiscoverer{
		messenger: messenger,
	}
}

func (d *P2pDiscoverer) DiscoverNewTargets(existingTargets []Target) (targets []Target) {
	targets = existingTargets
	currentlyConnectedTargets := d.messenger.ConnectedAddresses()

	for idx, address := range currentlyConnectedTargets {
		targetAddress := strings.Split(address, "/")
		target := MakeTarget(uint(idx), targetAddress[1], targetAddress[2], targetAddress[4], "NEW")
		if !containsTarget(targets, target) {
			targets = append(targets, target)
		}
	}
	return
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *P2pDiscoverer) IsInterfaceNil() bool {
	return d == nil
}
