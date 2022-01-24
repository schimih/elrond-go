package evaluation

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type EvaluationReport struct {
	evaluatorFactory EvaluatorFactory
	EvaluatedTargets []EvaluatedTarget
}

func NewEvaluationReport(ef EvaluatorFactory) EvaluationReport {
	return EvaluationReport{
		evaluatorFactory: ef,
		EvaluatedTargets: make([]EvaluatedTarget, 0),
	}
}

func (eR *EvaluationReport) RunEvaluation(scanResults []scan.ScannedTarget, evaluationType utils.EvaluationType) (EvaluatedTargets []EvaluatedTarget) {
	switch evaluationType {
	case utils.PortStatusEvaluation:
		eR.evaluatePortStatus(scanResults)
	case utils.NoEvaluation:
		log.Error("No Evaluation Type given - Please choose evaluation type")
	default:
		log.Error("Unknown Evaluation Type")
	}
	return eR.EvaluatedTargets
}

func (eR *EvaluationReport) evaluatePortStatus(scanResults []scan.ScannedTarget) {
	for _, scannedTarget := range scanResults {
		if !find(scannedTarget.Address, eR.EvaluatedTargets) {
			eR.populateReport(scannedTarget)
		}
	}
}

func (eR *EvaluationReport) populateReport(scannedTarget scan.ScannedTarget) {
	evaluator := eR.evaluatorFactory.CreateEvaluator(scannedTarget.Address, scannedTarget.Ports, scannedTarget.AnalysisType)
	eR.EvaluatedTargets = append(eR.EvaluatedTargets, evaluator.Evaluate())
}

func (eR *EvaluationReport) GetNumberOfEvaluatedTargets() (length int) {
	return len(eR.EvaluatedTargets)
}

func (eR *EvaluationReport) GetEvaluationTargets() (EvaluatedTargets []EvaluatedTarget) {
	return eR.EvaluatedTargets
}

func find(needle string, haystack []EvaluatedTarget) bool {
	for _, EvaluationTarget := range haystack {
		if needle == EvaluationTarget.Address {
			return true
		}
	}
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (eR *EvaluationReport) IsInterfaceNil() bool {
	return eR == nil
}
