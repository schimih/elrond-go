package scan

import "github.com/elrond-go/cmd/vat/core"

type Peer struct {
	ID           uint
	Protocol     string
	Address      string
	Standard     string
	Ports        []Port
	Architecture string
	Status       core.TargetStatus
	AnalysisType int
}

func NewPeer(id uint, address string, ports []Port, targetStatus string, peerStatus core.TargetStatus, analysisType int) Peer {

	return Peer{
		ID:           id,
		Protocol:     "",
		Address:      address,
		Standard:     "",
		Ports:        ports,
		Architecture: "",
		Status:       peerStatus,
		AnalysisType: analysisType,
	}
}
