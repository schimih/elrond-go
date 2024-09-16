package export

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/elrond-go/cmd/vat/core"
	"github.com/stretchr/testify/require"
)

func TestXMLOutput(t *testing.T) {
	rankedReport := NewRankedReport()
	formatter := &XMLFormatter{}

	testTargetsSlice := CreateEvaluatedTargetsTestSlice(3, core.LOW, core.Closed, 3)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(2, core.MID, core.Closed, 5)...)
	testTargetsSlice = append(testTargetsSlice, CreateEvaluatedTargetsTestSlice(1, core.HIGH, core.Closed, 5)...)
	rankedReport.populateReport(testTargetsSlice)
	rankedReport.NodesAnalyzed = len(testTargetsSlice)
	rankedReport.sortReport()

	formatter.Output(rankedReport)
	src := "AnalysisResults.xml"
	JsonRankedReport, _ := openXMLFile(src, rankedReport)

	require.Equal(t, 6, len(JsonRankedReport.HighRiskNodes))
	require.Equal(t, 4, len(JsonRankedReport.MediumRiskNodes))
	require.Equal(t, 2, len(JsonRankedReport.LowRiskNodes))
}

func openXMLFile(src string, rankedReport RankedReport) (RankedReport, error) {
	fpath, _ := os.Stat(src)
	if filepath.Ext(fpath.Name()) != ".xml" {
		return rankedReport, fmt.Errorf("File not created.")
	}

	xmlFile, err := os.Open(src)
	if err != nil {
		return rankedReport, fmt.Errorf("Could not open File")
	}
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	//in peers we will have all results from xml file
	xml.Unmarshal(byteValue, &rankedReport)

	return rankedReport, nil
}
