package output

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type TableFormatter struct {
	OutputType       utils.OutputType
	EvaluationReport []evaluation.EvaluationTarget
}

// to be refactored
func (tF *TableFormatter) GetOutput() {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := tF.EvaluationReport
	if len(peersDB) == 0 {
		log.Info("No peers in DB. First load a json or run discovery!")
		return
	}
	dataLines := make([]*display.LineData, 0, len(tF.EvaluationReport))
	for idx, evaluationResult := range tF.EvaluationReport {
		rAddress := evaluationResult.Address
		if len(evaluationResult.Ports) != 0 {
			for jdx, tPort := range evaluationResult.Ports {
				if jdx == 0 {
					rPort := fmt.Sprintf("%d", tPort.Number)
					rStatus := string(utils.PortStatus(tPort.State))
					rProtocol := tPort.Protocol
					rIndex := fmt.Sprintf("%d", idx)
					horizontalLineAfter := jdx == len(evaluationResult.Ports)-1
					if len(evaluationResult.Ports) == 1 {
						horizontalLineAfter = false
					}
					lines := display.NewLineData(horizontalLineAfter, []string{rIndex, rAddress, rPort, rStatus, rProtocol})
					dataLines = append(dataLines, lines)
				} else {
					rPort := fmt.Sprintf("%d", tPort.Number)
					rStatus := string(utils.PortStatus(tPort.State))
					rProtocol := tPort.Protocol
					horizontalLineAfter := jdx == len(evaluationResult.Ports)-1
					lines := display.NewLineData(horizontalLineAfter, []string{"", "", rPort, rStatus, rProtocol})
					dataLines = append(dataLines, lines)
				}
			}
		} else {
			rIndex := fmt.Sprintf("%d", idx)
			lines := display.NewLineData(true, []string{rIndex, rAddress, "NO OPEN PORTS"})
			dataLines = append(dataLines, lines)
		}

		Score := fmt.Sprintf("%d", tF.EvaluationReport[idx].Score)
		//totalHorizontalLine := (temp == len(dataLines)+1)
		totalline := display.NewLineData(true, []string{">>>>>", ">>>>>>>>>>>", ">>>>>", "RATING", Score})
		dataLines = append(dataLines, totalline)
	}

	table, _ := display.CreateTableString(header, dataLines)
	fmt.Println(table)
}

// IsInterfaceNil returns true if there is no value under the interface
func (tF *TableFormatter) IsInterfaceNil() bool {
	return tF == nil
}
