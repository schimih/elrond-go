package export

import (
	"sort"

	"github.com/elrond-go/cmd/vat/evaluation"
	core "github.com/elrond-go/cmd/vat/utils"
)

type RankedReport struct {
	HighRiskNodes   []evaluation.EvaluatedTarget
	MediumRiskNodes []evaluation.EvaluatedTarget
	LowRiskNodes    []evaluation.EvaluatedTarget
	HiddenNodes     []evaluation.EvaluatedTarget
	NodesAnalyzed   int
}

func NewRankedReport() RankedReport {
	return RankedReport{
		HighRiskNodes:   make([]evaluation.EvaluatedTarget, 0),
		MediumRiskNodes: make([]evaluation.EvaluatedTarget, 0),
		LowRiskNodes:    make([]evaluation.EvaluatedTarget, 0),
		HiddenNodes:     make([]evaluation.EvaluatedTarget, 0),
		NodesAnalyzed:   0,
	}
}

func (rR *RankedReport) populateReport(evaluatedTargets []evaluation.EvaluatedTarget) {
	for _, evaluatedTarget := range evaluatedTargets {
		switch evaluatedTarget.GetSecurityLevel() {
		case core.HIGH:
			rR.LowRiskNodes = append(rR.LowRiskNodes, evaluatedTarget)
		case core.MID:
			rR.MediumRiskNodes = append(rR.MediumRiskNodes, evaluatedTarget)
		default:
			rR.HighRiskNodes = append(rR.HighRiskNodes, evaluatedTarget)
		}
	}
}

func (rR *RankedReport) sortReport() {
	// sort nodes from high risk to low risk
	sortSlice(rR.LowRiskNodes)

	sortSlice(rR.MediumRiskNodes)

	sortSlice(rR.HighRiskNodes)
}

func (rR *RankedReport) getAllEvaluatedTargets() (evaluatedTargets []evaluation.EvaluatedTarget) {

	evaluatedTargets = append(rR.HighRiskNodes, rR.MediumRiskNodes...) // can't concatenate more than 2 slice at once

	evaluatedTargets = append(evaluatedTargets, rR.LowRiskNodes...)

	return append(evaluatedTargets, rR.HiddenNodes...)
}

func sortSlice(slice []evaluation.EvaluatedTarget) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].GetScore() < slice[j].GetScore()
	})
}

// IsInterfaceNil returns true if there is no value under the interface
func (rR *RankedReport) IsInterfaceNil() bool {
	return rR == nil
}
