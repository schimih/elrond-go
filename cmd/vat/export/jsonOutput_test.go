package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/elrond-go/cmd/vat/utils"
	"github.com/stretchr/testify/require"
)

func TestJsonOutput(t *testing.T) {
	rankedReport := NewRankedReport()
	formatter := &JsonFormatter{}

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(3, utils.LOW, utils.Closed, 3)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(2, utils.MID, utils.Closed, 5)...)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, utils.HIGH, utils.Closed, 5)...)

	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	formatter.Output(rankedReport)
	src := "AnalysisResults.json"
	JsonRankedReport, _ := openJsonFile(src, rankedReport)

	require.Equal(t, 3, len(JsonRankedReport.HighRiskNodes))
	require.Equal(t, 2, len(JsonRankedReport.MediumRiskNodes))
	require.Equal(t, 1, len(JsonRankedReport.LowRiskNodes))
}

func openJsonFile(src string, rankedReport RankedReport) (RankedReport, error) {
	fpath, _ := os.Stat(src)
	if filepath.Ext(fpath.Name()) != ".json" {
		return rankedReport, fmt.Errorf("File not created.")
	}

	jsonFile, err := os.Open(src)
	if err != nil {
		return rankedReport, fmt.Errorf("Could not open File")
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//in peers we will have all results from json file
	json.Unmarshal(byteValue, &rankedReport)

	return rankedReport, nil
}
