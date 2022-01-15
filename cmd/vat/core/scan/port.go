package scan

import go_nmap "github.com/lair-framework/go-nmap"

type Port struct {
	ID       uint
	Number   int
	Protocol string
	State    string
	Owner    string
}

type Ports struct {
	Ports []Port
	Host  go_nmap.Host
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

func NewPort(id uint,
	number int,
	protocol string,
	state string,
	owner string) Port {

	return Port{
		ID:       id,
		Number:   number,
		Protocol: protocol,
		State:    state,
		Owner:    owner,
	}
}

func createPortSlice(host go_nmap.Host) Ports {
	return Ports{
		Ports: make([]Port, 0),
		Host:  host,
	}
}

func (ps *Ports) translatePortSlice() (portSlice []Port) {
	for idx, port := range ps.Host.Ports {
		ps.Ports = append(ps.Ports, NewPort(uint(idx), port.PortId, port.Protocol, port.State.State, port.Owner.Name))
	}
	return ps.Ports
}
