package result

import (
	go_nmap "github.com/lair-framework/go-nmap"
)

type Rating struct {
	State  string
	Reason []string
	Value  uint
}

type RatingState string

const (
	Assessed    RatingState = "assessed"
	Sweeped     RatingState = "sweeped"
	Unknown     RatingState = "unknown"
	NotAssessed RatingState = "not-assessed"
)

func (rating Rating) Status() RatingState {
	return RatingState(rating.State)
}

// PortStatus represents a port's status.
type PortStatus string

// Enumerates the different possible state values.
const (
	Open       PortStatus = "open"
	Closed     PortStatus = "closed"
	Filtered   PortStatus = "filtered"
	Unfiltered PortStatus = "unfiltered"
)

// Status returns the status of a port.
func (p Port) Status() PortStatus {
	return PortStatus(p.State)
}

func (r *ResultsContainer) Process(NmapScanResult []*go_nmap.NmapRun, testType string) *ResultsContainer {
	for _, nmapRun := range NmapScanResult {
		for hidx, host := range nmapRun.Hosts {
			portSlice := make([]Port, 0)
			for pidx, port := range host.Ports {
				port := NewPort(uint(pidx), port.PortId, port.Protocol, port.State.State, port.Owner.Name)
				portSlice = append(portSlice, port)
			}
			peer := NewPeer(uint(hidx), host.Addresses[0].Addr, portSlice, host.Status.State, "not-assessed", testType)
			r.Results = append(r.Results, peer)
		}
	}
	return r
}
