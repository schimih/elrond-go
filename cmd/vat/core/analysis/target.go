package analysis

import "github.com/elrond-go/cmd/vat/core/scan"

type Target struct {
	ID             uint
	Protocol       string
	Address        string
	ConnectionPort string
	Status         scan.TargetStatus
}

func MakeTarget(id uint, protocol string, address string, connectionPort string, status scan.TargetStatus) Target {
	return Target{
		ID:             id,
		Protocol:       protocol,
		Address:        address,
		ConnectionPort: connectionPort,
		Status:         status,
	}
}

func containsTarget(haystack []Target, needle Target) bool {
	for _, target := range haystack {
		if target.Address == needle.Address {
			return true
		}
	}
	return false
}
