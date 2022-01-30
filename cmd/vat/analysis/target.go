package analysis

import (
	"github.com/elrond-go/cmd/vat/utils"
)

// TargetStatus represents a target's state.
type TargetStatus string

// Enumerates the different possible state values.
const (
	New     TargetStatus = "NEW"
	Scanned TargetStatus = "SCANNED"
	Expired TargetStatus = "EXPIRED"
)

const (
	NoCommand = iota
	CHANGE_STATUS_TO_EXPIRED
)

// ActualStatus returns the status of a target.
func (t DiscoveredTarget) ActualStatus() TargetStatus {
	return TargetStatus(t.Status)
}

type DiscoveredTarget struct {
	ID             uint
	Protocol       string
	Address        string
	ConnectionPort string
	Status         utils.TargetStatus
}

func MakeTarget(id uint, protocol string, address string, connectionPort string, status utils.TargetStatus) DiscoveredTarget {
	return DiscoveredTarget{
		ID:             id,
		Protocol:       protocol,
		Address:        address,
		ConnectionPort: connectionPort,
		Status:         status,
	}
}

func containsTarget(haystack []DiscoveredTarget, needle DiscoveredTarget) bool {
	for _, target := range haystack {
		if target.Address == needle.Address {
			return true
		}
	}
	return false
}
