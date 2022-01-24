package scan

import (
	"github.com/elrond-go/cmd/vat/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type Port struct {
	ID         uint
	Number     int
	Protocol   string
	State      utils.PortStatus
	Owner      string
	Type       utils.PortType
	RiskValue  utils.Risk
	RiskReason utils.Judgement
}

type Ports struct {
	Ports []Port
	Host  go_nmap.Host
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
		RiskValue:  100,
		RiskReason: utils.JudgementNoRisk,
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
		newPort := NewPort(uint(idx), port.PortId, port.Protocol, utils.PortStatus(port.State.State), port.Owner.Name, utils.Unknown)
		newPort.depictTypeAndImportance()
		ps.Ports = append(ps.Ports, newPort)
	}
	return ps.Ports
}

func (p *Port) depictTypeAndImportance() {
	if p.isPortInRange(37373, 38383) {
		p.Type = utils.ElrondPort
		p.RiskValue = utils.NoRisk
		p.RiskReason = utils.JudgementNoRisk
	} else if (p.Number == 80) || (p.Number == 8080) {
		p.Type = utils.WebPort
		p.RiskValue = utils.SmallRisk
		p.RiskReason = utils.JudgementWeb
	} else if p.Number == 22 {
		p.Type = utils.SshPort
		p.RiskValue = utils.SmallRisk
		p.RiskReason = utils.JudgementSsh
	} else {
		p.Type = utils.OutsideElrond
		p.RiskValue = utils.MediumRisk
		p.RiskReason = utils.JudgementMediumRisk
	}
}

func (p *Port) isPortInRange(low int, high int) bool {
	if (p.Number > low) && (p.Number < high) {
		return true
	}
	return false
}

// Status returns the status of a port.
func (p Port) Status() utils.PortStatus {
	return utils.PortStatus(p.State)
}
