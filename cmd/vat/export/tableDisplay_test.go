package export

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/elrond-go/cmd/vat/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableFormatterOutput(t *testing.T) {
	rankedReport := NewRankedReport()
	formatter := &TableFormatter{
		header:    make([]string, 0),
		dataLines: make([]*display.LineData, 0)}

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(1, core.LOW, core.Open, 5)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, core.MID, core.Open, 5)...)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, core.HIGH, core.Open, 5)...)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	formatter.Output(rankedReport)

	require.Equal(t, 5, len(formatter.header))     // 5 strings: -"Index", "Address", "Port", "Status", "Service"
	require.Equal(t, 21, len(formatter.dataLines)) // 5lines per target - 1 line per each port + 1 line for rating => 18 datalines
}

func TestTableFormatterOutput_NoEvaluatedTargetAvailable(t *testing.T) {
	rankedReport := NewRankedReport()
	formatter := &TableFormatter{
		header:    make([]string, 0),
		dataLines: make([]*display.LineData, 0)}
	rankedReport = RankedReport{}
	err := formatter.Output(rankedReport)

	expectedErrorString := "no scanned targets: nothing to display"
	assert.EqualErrorf(t, err, expectedErrorString, "wrong message")
}

func TestTableFormatterOutput_NoAccessiblePort(t *testing.T) {
	rankedReport := NewRankedReport()
	formatter := &TableFormatter{
		header:    make([]string, 0),
		dataLines: make([]*display.LineData, 0)}

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(1, core.LOW, core.Closed, 0)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, core.MID, core.Closed, 5)...)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, core.HIGH, core.Closed, 5)...)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	formatter.Output(rankedReport)

	require.Equal(t, 5, len(formatter.header))     // 5 strings: -"Index", "Address", "Port", "Status", "Service"
	require.Equal(t, 24, len(formatter.dataLines)) // 5lines per target (one is empty -> only no accessible port line) - 1 line per each port + 1 line for rating => 14 datalines
	require.Equal(t, "NO ACCESSIBLE PORTS", formatter.dataLines[0].Values[2])
}
