package manager

import (
	"time"

	"github.com/elrond-go/cmd/vat/core/evaluation"
)

var minutesInADay = 1440

type AnalysisManager struct {
	Start               time.Time
	TimeUntilExpiration time.Duration
	TotalRunTime        time.Duration
	AnalysisLoops       int
	NodesAnalyzed       int
	HighRiskNodes       []evaluation.Node
	MediumRiskNodes     []evaluation.Node
	LowRiskNodes        []evaluation.Node
	AnalysisType        int
	IsHidden            []evaluation.Node
	ExpireNextRun       bool
}

func NewAnalysisManager(analysisType int) AnalysisManager {
	return AnalysisManager{
		Start:               time.Now(),
		TimeUntilExpiration: time.Duration(0),
		TotalRunTime:        time.Duration(0),
		AnalysisLoops:       0,
		NodesAnalyzed:       0,
		HighRiskNodes:       make([]evaluation.Node, 0),
		MediumRiskNodes:     make([]evaluation.Node, 0),
		LowRiskNodes:        make([]evaluation.Node, 0),
		AnalysisType:        analysisType,
		IsHidden:            make([]evaluation.Node, 0),
		ExpireNextRun:       false,
	}
}
