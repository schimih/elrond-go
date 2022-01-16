package evaluation

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type EvaluationReport struct {
	evaluatorFactory  EvaluatorFactory
	EvaluationTargets []EvaluationTarget
	evaluationType    utils.EvaluationType
}

func NewEvaluationReport(ef EvaluatorFactory) EvaluationReport {
	return EvaluationReport{
		evaluatorFactory:  ef,
		EvaluationTargets: make([]EvaluationTarget, 0),
		evaluationType:    utils.NoEvaluation,
	}
}

func (eR *EvaluationReport) StartEvaluation(scanResults []scan.AnalyzedTarget, evaluationType utils.EvaluationType) {
	eR.evaluationType = evaluationType
	switch evaluationType {
	case utils.PortStatusEvaluation:
		eR.evaluatePortStatus(scanResults)
	case utils.NoEvaluation:
		log.Error("No Evaluation Type given - Please choose evaluation type")
	default:
		log.Error("Unknown Evaluation Type")
	}
}

func (eR *EvaluationReport) evaluatePortStatus(scanResults []scan.AnalyzedTarget) {
	for _, analyzedTarget := range scanResults {
		if !find(analyzedTarget.Address, eR.EvaluationTargets) {
			eR.populateReport(analyzedTarget)
		}
	}
}

func (eR *EvaluationReport) populateReport(peer scan.AnalyzedTarget) {
	evaluator := eR.evaluatorFactory.CreateEvaluator(peer.Address, peer.Ports, peer.AnalysisType)
	eR.EvaluationTargets = append(eR.EvaluationTargets, evaluator.Evaluate())
}

// to be refactored
func (eR *EvaluationReport) DisplayToTable() {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := eR.EvaluationTargets
	if len(peersDB) == 0 {
		log.Info("No peers in DB. First load a json or run discovery!")
		return
	}
	dataLines := make([]*display.LineData, 0, len(eR.EvaluationTargets))
	for idx, evaluationResult := range eR.EvaluationTargets {
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

		Score := fmt.Sprintf("%d", eR.EvaluationTargets[idx].Score)
		//totalHorizontalLine := (temp == len(dataLines)+1)
		totalline := display.NewLineData(true, []string{">>>>>", ">>>>>>>>>>>", ">>>>>", "RATING", Score})
		dataLines = append(dataLines, totalline)
	}

	table, _ := display.CreateTableString(header, dataLines)
	fmt.Println(table)
}

func find(needle string, haystack []EvaluationTarget) bool {
	for _, EvaluationTarget := range haystack {
		if needle == EvaluationTarget.Address {
			return true
		}
	}
	return false
}
