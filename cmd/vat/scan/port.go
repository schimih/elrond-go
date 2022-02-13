package scan

import (
	"github.com/elrond-go/cmd/vat/utils"
	gonmap "github.com/lair-framework/go-nmap"
)

type Port struct {
	ID         uint
	Number     int
	Protocol   string
	State      utils.PortStatus
	Owner      string
	Type       utils.PortType
	Importance utils.Judgement
}

type Ports struct {
	Ports []Port
	Host  gonmap.Host
}

func NewPort(id uint,
	number int,
	protocol string,
	state utils.PortStatus,
	owner string, portType utils.PortType) Port {

	return Port{
		ID:         id,
		Number:     number,
		Protocol:   protocol,
		State:      state,
		Owner:      owner,
		Type:       portType,
		Importance: utils.JudgementNoRisk,
	}
}

func createPortSlice(host gonmap.Host) Ports {
	return Ports{
		Ports: make([]Port, 0),
		Host:  host,
	}
}

func (ps *Ports) translatePortSlice() (portSlice []Port) {
	for idx, port := range ps.Host.Ports {
		newPort := NewPort(uint(idx), port.PortId, port.Protocol, utils.PortStatus(port.State.State), port.Owner.Name, utils.Unknown)
		_ = newPort.depictTypeAndImportance()
		ps.Ports = append(ps.Ports, newPort)
	}
	return ps.Ports
}

func (p *Port) depictTypeAndImportance() bool {
	if p.isPortInRange(37373, 38383) {
		p.Type = utils.ElrondPort
		p.Importance = utils.JudgementNoRisk

		return true
	}

	if (p.Number == 80) || (p.Number == 8080) {
		p.Type = utils.WebPort
		p.Importance = utils.JudgementWeb

		return true
	}

	if p.Number == 22 {
		p.Type = utils.SshPort
		p.Importance = utils.JudgementSsh

		return true
	}

	p.Type = utils.OutsideElrond
	p.Importance = utils.JudgementMediumRisk

	return true
}

func (p *Port) isPortInRange(low int, high int) bool {
	if (p.Number > low) && (p.Number < high) {
		return true
	}
	return false
}

// Status returns the status of a port.
func (p Port) Status() utils.PortStatus {
	return p.State
}
