package evaluation

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/analysis"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type Report struct {
	evaluatorFactory EvaluatorFactory
	EvaluatedTargets []EvaluatedTarget
	scannerFactory   analysis.ScannerFactory
	evaluationType   utils.EvaluationType
}

func NewEvaluationReport(eF EvaluatorFactory, sF analysis.ScannerFactory) Report {
	return Report{
		evaluatorFactory: eF,
		EvaluatedTargets: make([]EvaluatedTarget, 0),
		scannerFactory:   sF,
		evaluationType:   utils.NoEvaluation,
	}
}

func (eR *Report) RunEvaluation(scanResults []scan.ScannedTarget, evaluationType utils.EvaluationType) (EvaluatedTargets []EvaluatedTarget, err error) {
	eR.evaluationType = evaluationType
	if eR.evaluationType != utils.NoEvaluation {
		eR.runEvaluation(scanResults)
	} else {
		log.Error("No Evaluation Type given - Please choose evaluation type")
		return nil, utils.ErrNoEvaluationType
	}

	return eR.EvaluatedTargets, nil
}

func (eR *Report) runEvaluation(scanResults []scan.ScannedTarget) {
	for _, scannedTarget := range scanResults {
		if !find(scannedTarget.Address, eR.EvaluatedTargets) {
			eR.populateReport(scannedTarget)
		}
	}
}

func (eR *Report) populateReport(scannedTarget scan.ScannedTarget) {
	evaluator := eR.evaluatorFactory.CreateEvaluator(scannedTarget.Address, scannedTarget.Ports, eR.evaluationType, eR.scannerFactory)
	eR.EvaluatedTargets = append(eR.EvaluatedTargets, evaluator.Evaluate())
}

func find(needle string, haystack []EvaluatedTarget) bool {
	for _, evaluationTarget := range haystack {
		if needle == evaluationTarget.Address {
			return true
		}
	}
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (eR *Report) IsInterfaceNil() bool {
	return eR == nil
}
