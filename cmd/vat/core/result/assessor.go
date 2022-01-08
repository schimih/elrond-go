package result

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/display"
)

type assessor struct {
	assessments []AssessmentRunner
}

func (r *ResultsContainer) EvaluateNewPeers(container *ResultsContainer) {
	for _, peer := range r.Results {
		peer.TestType = "TCP-ELROND"
		if peer.TestType == "TCP-ELROND" {
			if peer.Evaluation.State != "EVALUATED" {
				r.RunAssessment()
			} else {
				log.Info("PEER already evaluated")
			}
		} else {
			log.Info("Please run TCP-ALL testType for full Evaluation")
		}
	}
}

func (r *ResultsContainer) RunAssessment() {
	for idx, result := range r.Results {
		for _, port := range result.Ports {
			if port.Status() == Open {
				if port.InsidePortRange(37373, 38383) {
					r.Results[idx].Evaluation.Process("TCP-ELROND")
				} else {
					if port.Number == 80 {
						r.Results[idx].Evaluation.Process("TCP-WEB")
					}
					if port.Number == 8080 {
						r.Results[idx].Evaluation.Process("TCP-WEB")
					}
					if port.Number == 22 {
						r.Results[idx].Evaluation.Process("TCP-SSH")

					}
				}
			} else {
				r.Results[idx].Evaluation.Process("TCP-SSH")
			}
		}

	}
}

func (r *ResultsContainer) DisplayAnalysisInfo() {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := r.Results
	if len(peersDB) == 0 {
		log.Info("No peers in DB. First load a json or run discovery!")
		return
	}
	dataLines := make([]*display.LineData, 0, len(r.Results))
	for idx, p := range r.Results {
		rAddress := p.Address
		if len(p.Ports) != 0 {
			for jdx, tPort := range p.Ports {
				if jdx == 0 {
					rPort := fmt.Sprintf("%d", tPort.Number)
					rStatus := tPort.State
					rProtocol := tPort.Protocol
					rIndex := fmt.Sprintf("%d", idx)
					horizontalLineAfter := jdx == len(p.Ports)-1
					if len(p.Ports) == 1 {
						horizontalLineAfter = false
					}
					lines := display.NewLineData(horizontalLineAfter, []string{rIndex, rAddress, rPort, rStatus, rProtocol})
					dataLines = append(dataLines, lines)
				}
				rPort := fmt.Sprintf("%d", tPort.Number)
				rStatus := tPort.State
				rProtocol := tPort.Protocol
				horizontalLineAfter := jdx == len(p.Ports)-1
				lines := display.NewLineData(horizontalLineAfter, []string{"", "", rPort, rStatus, rProtocol})
				dataLines = append(dataLines, lines)
			}
		} else {
			rIndex := fmt.Sprintf("%d", idx)
			lines := display.NewLineData(true, []string{rIndex, rAddress, "NO OPEN PORTS"})
			dataLines = append(dataLines, lines)
		}

		Rating := fmt.Sprintf("%d", r.Results[idx].Evaluation.Value)
		//totalHorizontalLine := (temp == len(dataLines)+1)
		totalline := display.NewLineData(true, []string{">>>>>", ">>>>>>>>>>>", ">>>>>", "RATING", Rating})
		dataLines = append(dataLines, totalline)
	}

	table, _ := display.CreateTableString(header, dataLines)
	fmt.Println(table)
}
