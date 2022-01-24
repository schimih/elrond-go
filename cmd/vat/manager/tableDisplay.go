package manager

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/utils"
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
		return fmt.Errorf("No scanned targets: nothing to display.")
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
					string(utils.PortStatus(tPort.State)),
					tPort.Protocol})
			} else {
				line = display.NewLineData(horizontalLineAfter, []string{
					"",
					"",
					fmt.Sprintf("%d", tPort.Number),
					string(utils.PortStatus(tPort.State)),
					tPort.Protocol})
			}
			tF.dataLines = append(tF.dataLines, line)
		}
	} else {
		line = display.NewLineData(true, []string{
			fmt.Sprintf("%d", id),
			evaluationResult.Address,
			"NO ACCESIBLE PORTS"})
		tF.dataLines = append(tF.dataLines, line)
	}
	tF.addRating(evaluationResult)
}

func (tF *TableFormatter) addRating(evaluationResult evaluation.EvaluatedTarget) {
	totalline := display.NewLineData(true, []string{">>>>>",
		">>>>>>>>>>>",
		">>>>>",
		"RATING",
		fmt.Sprintf("%d", evaluationResult.Score)})
	tF.dataLines = append(tF.dataLines, totalline)
}

// IsInterfaceNil returns true if there is no value under the interface
func (tF *TableFormatter) IsInterfaceNil() bool {
	return tF == nil
}
