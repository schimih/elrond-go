package evaluation

import "testing"

func TestNewEvaluationReport(t *testing.T) {
	eF := NewEvaluatorFactory()

	eR := NewEvaluationReport(eF)
}
