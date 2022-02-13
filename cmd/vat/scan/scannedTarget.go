package scan

import (
	core "github.com/elrond-go/cmd/vat/core"
)

type ScannedTarget struct {
	ID           uint
	Protocol     string
	Address      string
	Standard     string
	Ports        []Port
	Architecture string
	Status       core.TargetStatus
	AnalysisType core.AnalysisType
}

func NewScannedTarget(id uint, address string, ports []Port, peerStatus core.TargetStatus, analysisType core.AnalysisType) ScannedTarget {

	return ScannedTarget{
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
