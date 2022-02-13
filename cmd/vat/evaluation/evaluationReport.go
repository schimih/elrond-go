package evaluation

import (
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/analysis"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type report struct {
	mut              sync.Mutex
	evaluatorFactory EvaluatorFactory
	evaluatedTargets []EvaluatedTarget
	scannerFactory   analysis.ScannerFactory
	evaluationType   utils.EvaluationType
}

func NewReport(eF EvaluatorFactory, sF analysis.ScannerFactory) *report {
	return &report{
		evaluatorFactory: eF,
		evaluatedTargets: make([]EvaluatedTarget, 0),
		scannerFactory:   sF,
		evaluationType:   utils.NoEvaluation,
	}
}

func (r *report) RunEvaluation(scanResults []scan.ScannedTarget, evaluationType utils.EvaluationType) (evaluatedTargets []EvaluatedTarget, err error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.evaluationType = evaluationType
	if r.evaluationType == utils.NoEvaluation {
		log.Error("No Evaluation Type given - Please choose evaluation type")
		return nil, utils.ErrNoEvaluationType
	}

	r.runEvaluation(scanResults)
	return r.evaluatedTargets, nil
}

func (r *report) runEvaluation(scanResults []scan.ScannedTarget) {
	for _, scannedTarget := range scanResults {
		if !find(scannedTarget.Address, r.evaluatedTargets) {
			r.populateReport(scannedTarget)
		}
	}
}

func (r *report) populateReport(scannedTarget scan.ScannedTarget) {
	evaluator := r.evaluatorFactory.CreateEvaluator(scannedTarget.Address, scannedTarget.Ports, r.evaluationType, r.scannerFactory)
	r.evaluatedTargets = append(r.evaluatedTargets, evaluator.Evaluate())
}

func find(needle string, haystack []EvaluatedTarget) bool {
	for _, evaluationTarget := range haystack {
		if needle == evaluationTarget.identity.address {
			return true
		}
	}
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (r *report) IsInterfaceNil() bool {
	return r == nil
}
