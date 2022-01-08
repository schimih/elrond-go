package test

import (
	"reflect"
	"strings"
)

type Target struct {
	ID             uint
	Protocol       string
	Address        string
	ConnectionPort string
	Status         string
}

func (t *Test) discoverPeers() (discoveredPeers []Target) {
	discoveredPeers = make([]Target, 0)
	var target Target

	for idx, address := range t.Messenger.ConnectedAddresses() {
		peerAddress := strings.Split(address, "/")
		discoveredPeers = append(discoveredPeers, target.New(uint(idx), peerAddress[1], peerAddress[2], peerAddress[4], "NEW"))
	}
	return discoveredPeers
}

func deduplicateDiscoveredPeersList(newlyDiscoveredPeers []Target, TargetsList []Target) (updatedTargetsList []Target) {
	updatedTargetsList = TargetsList
	for _, target := range newlyDiscoveredPeers {
		exists := false
		for j, _ := range TargetsList {
			exists = reflect.DeepEqual(target.Address, TargetsList[j].Address)
			if exists {
				break
			}
		}
		if !exists {
			updatedTargetsList = append(updatedTargetsList, target)
			log.Info("Added target with address", "address", target.Address)
		}
	}
	return updatedTargetsList
}

func (t *Test) UpdateTargetsList() {
	t.TargetsList = deduplicateDiscoveredPeersList(t.discoverPeers(), t.TargetsList)
}

func (target *Target) New(id uint, protocol string, address string, connectionPort string, status string) Target {
	target = &Target{
		ID:             id,
		Protocol:       protocol,
		Address:        address,
		ConnectionPort: connectionPort,
		Status:         status,
	}
	return *target
}
