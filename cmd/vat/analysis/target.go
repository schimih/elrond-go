package analysis

import (
	"github.com/elrond-go/cmd/vat/core"
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
	NO_COMMAND = iota
	CHANGE_STATUS_TO_EXPIRED
)

// Status returns the status of a target.
func (t Target) ActualStatus() TargetStatus {
	return TargetStatus(t.Status)
}

type Target struct {
	ID             uint
	Protocol       string
	Address        string
	ConnectionPort string
	Status         core.TargetStatus
}

func MakeTarget(id uint, protocol string, address string, connectionPort string, status core.TargetStatus) Target {
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
