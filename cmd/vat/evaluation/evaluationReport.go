package evaluation

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/analysis"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type EvaluationReport struct {
	evaluatorFactory EvaluatorFactory
	EvaluatedTargets []EvaluatedTarget
	scannerFactory   analysis.ScannerFactory
	evaluationType   utils.EvaluationType
}

func NewEvaluationReport(eF EvaluatorFactory, sF analysis.ScannerFactory) EvaluationReport {
	return EvaluationReport{
		evaluatorFactory: eF,
		EvaluatedTargets: make([]EvaluatedTarget, 0),
		scannerFactory:   sF,
		evaluationType:   utils.NoEvaluation,
	}
}

func (eR *EvaluationReport) RunEvaluation(scanResults []scan.ScannedTarget, evaluationType utils.EvaluationType) (EvaluatedTargets []EvaluatedTarget) {
	eR.evaluationType = evaluationType
	switch eR.evaluationType {
	case utils.PortStatusEvaluation:
		eR.prepareEvaluation(scanResults)
	case utils.Polite_PortAndSshEvaluation:
		eR.prepareEvaluation(scanResults)
	case utils.NoEvaluation:
		log.Error("No Evaluation Type given - Please choose evaluation type")
	default:
		log.Error("Unknown Evaluation Type")
	}
	return eR.EvaluatedTargets
}

func (eR *EvaluationReport) prepareEvaluation(scanResults []scan.ScannedTarget) {
	for _, scannedTarget := range scanResults {
		if !find(scannedTarget.Address, eR.EvaluatedTargets) {
			eR.populateReport(scannedTarget)
		}
	}
}

func (eR *EvaluationReport) populateReport(scannedTarget scan.ScannedTarget) {
	evaluator := eR.evaluatorFactory.CreateEvaluator(scannedTarget.Address, scannedTarget.Ports, eR.evaluationType, eR.scannerFactory)
	eR.EvaluatedTargets = append(eR.EvaluatedTargets, evaluator.Evaluate())
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
