package export

import (
	"testing"

	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestNewRankedReport(t *testing.T) {
	rankedReport := NewRankedReport()

	require.False(t, rankedReport.IsInterfaceNil())
}

func TestRankedReport_SortAndPopulateOneLowRiskTarget(t *testing.T) {
	rankedReport := NewRankedReport()

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(1, utils.HIGH, utils.Open, 5)

	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	require.Equal(t, 1, len(rankedReport.LowRiskNodes))
	require.Equal(t, 0, len(rankedReport.MediumRiskNodes))
	require.Equal(t, 0, len(rankedReport.HighRiskNodes))
}

func TestRankedReport_SortAndPopulateTwoLowRiskTargets(t *testing.T) {
	rankedReport := NewRankedReport()

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(2, utils.HIGH, utils.Open, 5)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	require.Equal(t, 2, len(rankedReport.LowRiskNodes))
	require.Equal(t, 0, len(rankedReport.MediumRiskNodes))
	require.Equal(t, 0, len(rankedReport.HighRiskNodes))
}

func TestRankedReport_SortAndPopulateTwoMediumRiskTargets(t *testing.T) {
	rankedReport := NewRankedReport()

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(2, utils.MID, utils.Open, 5)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	require.Equal(t, 0, len(rankedReport.LowRiskNodes))
	require.Equal(t, 2, len(rankedReport.MediumRiskNodes))
	require.Equal(t, 0, len(rankedReport.HighRiskNodes))
}

func TestRankedReport_SortAndPopulateTwoHighRiskTargets(t *testing.T) {
	rankedReport := NewRankedReport()

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(2, utils.LOW, utils.Open, 5)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	require.Equal(t, 0, len(rankedReport.LowRiskNodes))
	require.Equal(t, 0, len(rankedReport.MediumRiskNodes))
	require.Equal(t, 2, len(rankedReport.HighRiskNodes))
}

func TestRankedReport_GetAllEvaluatedTargets(t *testing.T) {
	rankedReport := NewRankedReport()

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(1, utils.LOW, utils.Open, 5)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, utils.MID, utils.Open, 5)...)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, utils.HIGH, utils.Open, 5)...)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	allEvaluatedTargetsFromTestReport := rankedReport.getAllEvaluatedTargets()

	require.Equal(t, 3, len(allEvaluatedTargetsFromTestReport))
}

func CreateEvaluatedTargetsTestSlice(targets int, risk utils.SecureLevel, portStatus utils.PortStatus, noPorts int) (evaluatedTargets []evaluation.EvaluatedTarget) {
	ports := make([]scan.Port, 0)
	judgements := make([]string, 0)
	testTargets := make([]evaluation.EvaluatedTarget, 0)
	for i := 0; i < noPorts; i++ {
		ports = append(ports, scan.NewPort(0,
			0,
			"test",
			portStatus,
			"test", utils.PortType(utils.TCP_ELROND)))
		judgements = append(judgements, "test Judgement")
	}

	target := evaluation.EvaluatedTarget{
		//identity:   newIdentity("test", ports),
		//evaluation: newEvaluationResult(utils.NoEvaluation),
	}
	for i := 0; i < targets; i++ {
		testTargets = append(testTargets, target)
	}

	return testTargets
}
