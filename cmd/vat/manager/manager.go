package manager

import (
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/utils"
)

var minutesInADay = 1440

type AnalysisManager struct {
	Start               time.Time
	TimeUntilExpiration time.Duration
	TotalRunTime        time.Duration
	AnalysisLoops       int
	RankedReport        RankedReport
	FormatterFactory    *FormatterFactory
	formatterType       utils.OutputType
	AnalysisType        utils.AnalysisType
	EvaluationType      utils.EvaluationType
	FormatType          utils.OutputType
	ExpireNextRun       bool
}

func NewAnalysisManager(fF *FormatterFactory) (*AnalysisManager, error) {

	if check.IfNil(fF) {
		return nil, fmt.Errorf("FormatterFactory needed")
	}
	rankedReport := NewRankedReport()
	return &AnalysisManager{
		Start:               time.Now(),
		TimeUntilExpiration: time.Duration(0),
		TotalRunTime:        time.Duration(0),
		AnalysisLoops:       0,
		RankedReport:        rankedReport,
		FormatterFactory:    fF,
		AnalysisType:        utils.TCP_WEB,              // by default go with TCP_WEB -> this has to be controled by manager
		EvaluationType:      utils.PortStatusEvaluation, // by default go with PortStatusEvaluation
		FormatType:          utils.JSON,
		ExpireNextRun:       false,
	}, nil
}

// CompleteRound shall make sure that activities are completed efficiently and according to commands:
// - tbd - keep the timing
//
// - tbd - decide when discovered targets are expired so that they should be re-analyzed
//
// - plan the charactestics of the next round
//
// - create and control the formatter -
//
// - evaluate the status of the entire elrond ecosystem
//
// - organize results ... etc etc.
func (aM *AnalysisManager) CompleteRound(evaluatedTargets []evaluation.EvaluatedTarget) {

	aM.RankedReport.SortAndPopulate(evaluatedTargets)

	formatter := aM.FormatterFactory.CreateFormatter(aM.FormatType)
	formatter.Output(aM.RankedReport)

	aM.AnalysisLoops++
	aM.RankedReport = RankedReport{}
}

// IsInterfaceNil returns true if there is no value under the interface
func (aM *AnalysisManager) IsInterfaceNil() bool {
	return aM == nil
}
