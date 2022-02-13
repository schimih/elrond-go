package evaluation

import (
	"testing"

	"github.com/elrond-go/cmd/vat/scan/factory"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestNewEvaluatorFactory(t *testing.T) {
	NewFakeEvaluatorFactory := NewEvaluatorFactory()
	FakeEvaluatorFactory := EvaluatorFactory{}

	require.Equal(t, NewFakeEvaluatorFactory, FakeEvaluatorFactory)
}

func TestCreateEvaluator(t *testing.T) {
	fakeEvaluatedTarget := createFakeEvaluatedTarget(utils.PortStatusEvaluation, 20, 22, utils.SshPort)
	NewFakeEvaluatorFactory := NewEvaluatorFactory()
	sF := factory.NewScannerFactory()
	evaluator := NewFakeEvaluatorFactory.CreateEvaluator("testAddress", createFakePortsSlice(2, utils.Open, utils.SshPort, 2), utils.NoEvaluation, sF)
	evaluatedTarget := evaluator.Evaluate()

	require.Equal(t, fakeEvaluatedTarget.GetStatus(), evaluatedTarget.GetStatus())
}
