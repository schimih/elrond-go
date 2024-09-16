package evaluation

import (
	"testing"

	"github.com/elrond-go/cmd/vat/core"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/stretchr/testify/assert"
)

func TestNewEvaluationReport(t *testing.T) {
	evaluationReport := report{
		evaluatedTargets: make([]EvaluatedTarget, 0),
		scannerFactory:   scan.NewScannerFactory(),
		evaluationType:   core.NoEvaluation,
	}

	newEvaluationReport := NewReport(scan.NewScannerFactory())

	assert.Equal(t, evaluationReport, newEvaluationReport)
}

func TestRunEvaluation(t *testing.T) {
	newEvaluationReport := NewReport(scan.NewScannerFactory())

	someEmptyScannedResults := make([]scan.ScannedTarget, 1)
	evaluatedTargets, err := newEvaluationReport.RunEvaluation(someEmptyScannedResults, core.Polite_PortAndSshEvaluation)

	assert.Equal(t, 1, len(evaluatedTargets))
	assert.Nil(t, err)
}

func TestRunEvaluation_NoEvaluation(t *testing.T) {
	newEvaluationReport := NewReport(scan.NewScannerFactory())

	someEmptyScannedResults := make([]scan.ScannedTarget, 1)
	evaluatedTargets, err := newEvaluationReport.RunEvaluation(someEmptyScannedResults, core.NoEvaluation)

	assert.Equal(t, core.ErrNoEvaluationType, err)
	assert.Nil(t, evaluatedTargets)
}
