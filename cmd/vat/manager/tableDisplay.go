package manager

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/evaluation"
)

var log = logger.GetOrCreate("vat")

type TableFormatter struct {
	header    []string
	dataLines []*display.LineData
}

// to be refactored
func (tF *TableFormatter) Output(rankedReport RankedReport) error {
	if rankedReport.NodesAnalyzed == 0 {
		log.Info("No scanned targets: nothing to display.")
		return fmt.Errorf("no scanned targets: nothing to display")
	}

	tF.header = []string{"Index", "Address", "Port", "Status", "Service"}

	for idx, evaluationTarget := range rankedReport.GetAllEvaluatedTargets() {
		tF.addTargetToTable(idx, evaluationTarget)
	}

	table, _ := display.CreateTableString(tF.header, tF.dataLines)
	fmt.Println(table)

	return nil
}

func (tF *TableFormatter) addTargetToTable(id int, evaluationResult evaluation.EvaluatedTarget) {
	var line *display.LineData

	if len(evaluationResult.Ports) != 0 {
		for jdx, tPort := range evaluationResult.Ports {
			horizontalLineAfter := jdx == len(evaluationResult.Ports)-1
			if jdx == 0 {
				line = display.NewLineData(horizontalLineAfter, []string{
					fmt.Sprintf("%d", id),
					evaluationResult.Address,
					fmt.Sprintf("%d", tPort.Number),
					string(tPort.State),
					tPort.Protocol})
			} else {
				line = display.NewLineData(horizontalLineAfter, []string{
					"",
					"",
					fmt.Sprintf("%d", tPort.Number),
					string(tPort.State),
					tPort.Protocol})
			}
			tF.dataLines = append(tF.dataLines, line)
		}
	} else {
		line = display.NewLineData(true, []string{
			fmt.Sprintf("%d", id),
			evaluationResult.Address,
			"NO ACCESSIBLE PORTS"})
		tF.dataLines = append(tF.dataLines, line)
	}

	tF.addJudgement(evaluationResult)
	tF.addRating(evaluationResult)
}

func (tF *TableFormatter) addRating(evaluationResult evaluation.EvaluatedTarget) {
	totalLine := display.NewLineData(true, []string{">>>>>",
		">>>>>>>>>>>",
		">>>>>",
		"RATING",
		fmt.Sprintf("%d", evaluationResult.Score)})
	tF.dataLines = append(tF.dataLines, totalLine)
}

func (tF *TableFormatter) addJudgement(evaluationResult evaluation.EvaluatedTarget) {
	for _, judgement := range evaluationResult.Judgements {
		judgementLine := display.NewLineData(false, []string{"", "", "", "", judgement})
		tF.dataLines = append(tF.dataLines, judgementLine)
	}

}

// IsInterfaceNil returns true if there is no value under the interface
func (tF *TableFormatter) IsInterfaceNil() bool {
	return tF == nil
}
