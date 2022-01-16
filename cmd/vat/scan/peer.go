package scan

import (
	"github.com/elrond-go/cmd/vat/utils"
)

type Peer struct {
	ID           uint
	Protocol     string
	Address      string
	Standard     string
	Ports        []Port
	Architecture string
	Status       utils.TargetStatus
	AnalysisType utils.AnalysisType
}

func NewPeer(id uint, address string, ports []Port, targetStatus string, peerStatus utils.TargetStatus, analysisType utils.AnalysisType) Peer {

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
