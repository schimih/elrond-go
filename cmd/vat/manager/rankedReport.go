package manager

import (
	"sort"

	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/utils"
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

func (rR *RankedReport) SortAndPopulate(evaluatedTargets []evaluation.EvaluatedTarget) {
	for _, evaluatedTarget := range evaluatedTargets {
		switch evaluatedTarget.SecurityLevel {
		case utils.HIGH:
			rR.LowRiskNodes = append(rR.LowRiskNodes, evaluatedTarget)
		case utils.MID:
			rR.MediumRiskNodes = append(rR.MediumRiskNodes, evaluatedTarget)
		default:
			rR.HighRiskNodes = append(rR.HighRiskNodes, evaluatedTarget)
		}
		rR.NodesAnalyzed++
	}

	sort.Slice(rR.LowRiskNodes, func(i, j int) bool {
		return rR.LowRiskNodes[i].Score < rR.LowRiskNodes[j].Score
	})

	sort.Slice(rR.MediumRiskNodes, func(i, j int) bool {
		return rR.MediumRiskNodes[i].Score < rR.MediumRiskNodes[j].Score
	})

	sort.Slice(rR.HighRiskNodes, func(i, j int) bool {
		return rR.HighRiskNodes[i].Score < rR.HighRiskNodes[j].Score
	})
}

func (rR *RankedReport) GetAllEvaluatedTargets() (evaluatedTargets []evaluation.EvaluatedTarget) {

	for _, evaluatedTarget := range rR.HighRiskNodes {
		evaluatedTargets = append(evaluatedTargets, evaluatedTarget)
	}

	for _, evaluatedTarget := range rR.MediumRiskNodes {
		evaluatedTargets = append(evaluatedTargets, evaluatedTarget)
	}

	for _, evaluatedTarget := range rR.LowRiskNodes {
		evaluatedTargets = append(evaluatedTargets, evaluatedTarget)
	}

	for _, evaluatedTarget := range rR.HiddenNodes {
		evaluatedTargets = append(evaluatedTargets, evaluatedTarget)
	}

	return
}

// IsInterfaceNil returns true if there is no value under the interface
func (rR *RankedReport) IsInterfaceNil() bool {
	return rR == nil
}
