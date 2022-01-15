package evaluation

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/core/scan"
)

var log = logger.GetOrCreate("vat")

type AssessmentReport struct {
	evaluatorFactory  EvaluatorFactory
	EvaluationResults []EvaluationResult
}

func NewAssessmentReport(ef EvaluatorFactory) AssessmentReport {
	return AssessmentReport{
		evaluatorFactory:  ef,
		EvaluationResults: make([]EvaluationResult, 0),
	}
}

func (ar *AssessmentReport) GenerateReport(scanResults []scan.Peer) {
	for _, peer := range scanResults {
		if !find(peer.Address, ar.EvaluationResults) {
			ar.populateReport(peer)
		}
	}
}

func (ar *AssessmentReport) populateReport(peer scan.Peer) {
	node := NewNode(peer.Address, peer.Ports)
	evaluator := ar.evaluatorFactory.CreateEvaluator(node, peer.AnalysisType)
	ar.EvaluationResults = append(ar.EvaluationResults, evaluator.Evaluate(node))
}

// to be refactored
func (ar *AssessmentReport) DisplayToTable() {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := ar.EvaluationResults
	if len(peersDB) == 0 {
		log.Info("No peers in DB. First load a json or run discovery!")
		return
	}
	dataLines := make([]*display.LineData, 0, len(ar.EvaluationResults))
	for idx, p := range ar.EvaluationResults {
		rAddress := p.Node.Address
		if len(p.Node.Ports) != 0 {
			for jdx, tPort := range p.Node.Ports {
				if jdx == 0 {
					rPort := fmt.Sprintf("%d", tPort.Number)
					rStatus := tPort.State
					rProtocol := tPort.Protocol
					rIndex := fmt.Sprintf("%d", idx)
					horizontalLineAfter := jdx == len(p.Node.Ports)-1
					if len(p.Node.Ports) == 1 {
						horizontalLineAfter = false
					}
					lines := display.NewLineData(horizontalLineAfter, []string{rIndex, rAddress, rPort, rStatus, rProtocol})
					dataLines = append(dataLines, lines)
				} else {
					rPort := fmt.Sprintf("%d", tPort.Number)
					rStatus := tPort.State
					rProtocol := tPort.Protocol
					horizontalLineAfter := jdx == len(p.Node.Ports)-1
					lines := display.NewLineData(horizontalLineAfter, []string{"", "", rPort, rStatus, rProtocol})
					dataLines = append(dataLines, lines)
				}
			}
		} else {
			rIndex := fmt.Sprintf("%d", idx)
			lines := display.NewLineData(true, []string{rIndex, rAddress, "NO OPEN PORTS"})
			dataLines = append(dataLines, lines)
		}

		Score := fmt.Sprintf("%d", ar.EvaluationResults[idx].Score)
		//totalHorizontalLine := (temp == len(dataLines)+1)
		totalline := display.NewLineData(true, []string{">>>>>", ">>>>>>>>>>>", ">>>>>", "RATING", Score})
		dataLines = append(dataLines, totalline)
	}

	table, _ := display.CreateTableString(header, dataLines)
	fmt.Println(table)
}

func find(needle string, haystack []EvaluationResult) bool {
	for _, node := range haystack {
		if needle == node.Node.Address {
			return true
		}
	}
	return false
}
