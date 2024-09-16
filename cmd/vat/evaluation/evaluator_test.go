package evaluation

import (
	"testing"

	"github.com/elrond-go/cmd/vat/core"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/stretchr/testify/require"
)

func TestTargetEvaluation_StatusEvaluated(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(core.PortStatusEvaluation, 20, 10000, core.OutsideElrond)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, string(core.EVALUATED), fakeEvaluatedTarget.GetStatus())
}

func TestTargetEvaluation_SecurityLevelHIGH(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(core.PortStatusEvaluation, 20, 10000, core.OutsideElrond)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, core.HIGH, fakeEvaluatedTarget.GetSecurityLevel())
}

func TestTargetEvaluation_SecurityLevelMID(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(core.PortStatusEvaluation, 20, 22, core.SshPort)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, core.HIGH, fakeEvaluatedTarget.GetSecurityLevel())
}

// to be tested with real values to work
func TestTargetEvaluation_RunSpecialCheck(t *testing.T) {
	fakeTarget := createFakeEvaluatedTarget(core.Polite_PortAndSshEvaluation, 2, 22, core.SshPort)
	fakeEvaluatedTarget := fakeTarget.Evaluate()

	require.Equal(t, fakeTarget.GetPortsSlice()[0].Number, 22)
	require.Equal(t, 2, len(fakeTarget.GetPortsSlice()))
	require.Equal(t, fakeEvaluatedTarget.GetEvaluationType(), fakeTarget.GetEvaluationType())
	require.Equal(t, core.LOW, fakeEvaluatedTarget.GetSecurityLevel())
}

func createFakeEvaluatedTarget(evaluationType core.EvaluationType, noPorts int, portsNumber int, portType core.PortType) EvaluatedTarget {
	return EvaluatedTarget{
		identity:       newIdentity("168.119.106.29", createFakePortsSlice(portsNumber, core.Open, portType, noPorts)),
		evaluation:     newEvaluationResult(evaluationType),
		scannerFactory: scan.NewScannerFactory(),
	}
}

func createFakePortsSlice(number int, portStatus core.PortStatus, portType core.PortType, noPorts int) (portsSlice []scan.Port) {
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
