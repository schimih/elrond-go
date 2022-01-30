package manager

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/elrond-go/cmd/vat/utils"
)

type XMLFormatter struct {
}

func (xF *XMLFormatter) Output(rankedReport RankedReport) error {
	xmlData, _ := xml.MarshalIndent(rankedReport, "", " ")

	err := ioutil.WriteFile("AnalysisResults.xml", xmlData, 0644)
	if err != nil {
		return fmt.Errorf("could not write File")
	}
	path := utils.XMLFilePath
	log.Info("Evaluated Targets list added to ", "path", path)
	log.Info("Evaluated", "Nodes", rankedReport.NodesAnalyzed)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (xF *XMLFormatter) IsInterfaceNil() bool {
	return xF == nil
}
