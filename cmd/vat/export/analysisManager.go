package export

import (
	"time"

	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/elrond-go/cmd/vat/core"
	"github.com/elrond-go/cmd/vat/evaluation"
)

type AnalysisManager struct {
	Start               time.Time
	TimeUntilExpiration time.Duration
	TotalRunTime        time.Duration
	AnalysisLoops       int
	RankedReport        RankedReport
	formatterType       OutputType
	AnalysisType        core.AnalysisType
	EvaluationType      core.EvaluationType
	FormatType          OutputType
	ExpireNextRun       bool
}

func NewAnalysisManager() (*AnalysisManager, error) {

	rankedReport := NewRankedReport()
	return &AnalysisManager{
		Start:               time.Now(),
		TimeUntilExpiration: time.Duration(0),
		TotalRunTime:        time.Duration(0),
		AnalysisLoops:       0,
		RankedReport:        rankedReport,
		AnalysisType:        core.TCP_REQ1,                    // by default go with TCP_WEB -> this has to be controlled by manager
		EvaluationType:      core.Polite_PortAndSshEvaluation, // by default go with PortStatusEvaluation
		FormatType:          Table,
		ExpireNextRun:       false,
	}, nil
}

// CompleteRound shall make sure that activities are completed efficiently and according to commands:
// - tbd - keep the timing
//
// - tbd - decide when discovered targets are expired so that they should be re-analyzed
//
// - plan the characteristics of the next round
//
// - create and control the formatter -
//
// - evaluate the status of the entire elrond ecosystem
//
// - organize results ... etc. etc.
func (aM *AnalysisManager) CompleteRound(evaluatedTargets []evaluation.EvaluatedTarget) {
	aM.RankedReport.populateReport(evaluatedTargets)
	aM.RankedReport.NodesAnalyzed = len(evaluatedTargets)
	aM.RankedReport.sortReport()

	formatter, _ := createFormatter(aM.FormatType)
	err := formatter.Output(aM.RankedReport)
	if err != nil {
		return
	}

	aM.AnalysisLoops++
	aM.RankedReport = RankedReport{}
}

func createFormatter(formatType OutputType) (formatter Formatter, err error) {
	switch formatType {
	case Table:
		return &TableFormatter{
			header:    make([]string, 0),
			dataLines: make([]*display.LineData, 0),
		}, nil
	case JSON:
		return &JsonFormatter{}, nil
	case XML:
		return &XMLFormatter{}, nil
	default:
		return nil, core.ErrNoFormatterType
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (aM *AnalysisManager) IsInterfaceNil() bool {
	return aM == nil
}
