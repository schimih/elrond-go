package scan

import (
	"github.com/elrond-go/cmd/vat/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type Port struct {
	ID       uint
	Number   int
	Protocol string
	State    utils.PortStatus
	Owner    string
}

type Ports struct {
	Ports []Port
	Host  go_nmap.Host
}

func NewPort(id uint,
	number int,
	protocol string,
	state utils.PortStatus,
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
		ps.Ports = append(ps.Ports, NewPort(uint(idx), port.PortId, port.Protocol, utils.PortStatus(port.State.State), port.Owner.Name))
	}
	return ps.Ports
}

// Status returns the status of a port.
func (p Port) Status() utils.PortStatus {
	return utils.PortStatus(p.State)
}
