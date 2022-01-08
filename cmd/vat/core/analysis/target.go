package analysis

import "reflect"

type Target struct {
	ID             uint
	Protocol       string
	Address        string
	ConnectionPort string
	Status         string
}

func MakeTarget(id uint, protocol string, address string, connectionPort string, status string) Target {
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
		if reflect.DeepEqual(target.Address, needle.Address) {
			return true
		}
	}
	return false
}
