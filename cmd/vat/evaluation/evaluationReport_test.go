package evaluation

import (
	"testing"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/scan/factory"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewEvaluationReport(t *testing.T) {
	evaluationReport := report{
		evaluatorFactory: NewEvaluatorFactory(),
		evaluatedTargets: make([]EvaluatedTarget, 0),
		scannerFactory:   factory.NewScannerFactory(),
		evaluationType:   utils.NoEvaluation,
	}

	newEvaluationReport := NewReport(NewEvaluatorFactory(), factory.NewScannerFactory())

	assert.Equal(t, evaluationReport, newEvaluationReport)
}

func TestRunEvaluation(t *testing.T) {
	newEvaluationReport := NewReport(NewEvaluatorFactory(), factory.NewScannerFactory())

	someEmptyScannedResults := make([]scan.ScannedTarget, 1)
	evaluatedTargets, err := newEvaluationReport.RunEvaluation(someEmptyScannedResults, utils.Polite_PortAndSshEvaluation)

	assert.Equal(t, 1, len(evaluatedTargets))
	assert.Nil(t, err)
}

func TestRunEvaluation_NoEvaluation(t *testing.T) {
	newEvaluationReport := NewReport(NewEvaluatorFactory(), factory.NewScannerFactory())

	someEmptyScannedResults := make([]scan.ScannedTarget, 1)
	evaluatedTargets, err := newEvaluationReport.RunEvaluation(someEmptyScannedResults, utils.NoEvaluation)

	assert.Equal(t, utils.ErrNoEvaluationType, err)
	assert.Nil(t, evaluatedTargets)
}
