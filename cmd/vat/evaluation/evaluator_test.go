package evaluation

import (
	"testing"

	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/scan/factory"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestTargetEvaluation_StatusEvaluated(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(utils.PortStatusEvaluation, 20, 10000, utils.OutsideElrond)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, string(utils.EVALUATED), fakeEvaluatedTarget.GetStatus())
}

func TestTargetEvaluation_SecurityLevelHIGH(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(utils.PortStatusEvaluation, 20, 10000, utils.OutsideElrond)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, utils.HIGH, fakeEvaluatedTarget.GetSecurityLevel())
}

func TestTargetEvaluation_SecurityLevelMID(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(utils.PortStatusEvaluation, 20, 22, utils.SshPort)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, utils.HIGH, fakeEvaluatedTarget.GetSecurityLevel())
}

// to be tested with real values to work
func TestTargetEvaluation_RunSpecialCheck(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(utils.Polite_PortAndSshEvaluation, 2, 22, utils.SshPort)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, fakeTarget.GetPortsSlice()[0].Number, 22)
	require.Equal(t, 2, len(fakeTarget.GetPortsSlice()))
	require.Equal(t, fakeEvaluatedTarget.GetEvaluationType(), fakeTarget.GetEvaluationType())
	require.Equal(t, utils.LOW, fakeEvaluatedTarget.GetSecurityLevel())
}

func createFakeEvaluatedTarget(evaluationType utils.EvaluationType, noPorts int, portsNumber int, portType utils.PortType) EvaluatedTarget {
	return EvaluatedTarget{
		identity:       newIdentity("168.119.106.29", createFakePortsSlice(portsNumber, utils.Open, portType, noPorts)),
		evaluation:     newEvaluationResult(evaluationType),
		scannerFactory: factory.NewScannerFactory(),
	}
}

func createFakePortsSlice(number int, portStatus utils.PortStatus, portType utils.PortType, noPorts int) (portsSlice []scan.Port) {
	for i := 0; i < noPorts; i++ {
		portsSlice = append(portsSlice, scan.NewPort(0,
			number,
			"test",
			portStatus,
			"test", portType))
	}
	return
}

func createFakeJudgementsSlice(judgement string, nojudgements int) (judgmentsSlice []string) {
	for i := 0; i < nojudgements; i++ {
		judgmentsSlice = append(judgmentsSlice, judgement)
	}
	return
}
