package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elrond-go/cmd/vat/utils"
)

type JsonFormatter struct {
}

func (jF *JsonFormatter) Output(rankedReport RankedReport) error {
	jsonData, _ := json.MarshalIndent(rankedReport, "", " ")
	err := ioutil.WriteFile("AnalysisResults.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("could not write File")
	}

	path := utils.JsonFilePath
	log.Info("Evaluated Targets list added to ", "path", path)
	log.Info("Evaluated", "Nodes", rankedReport.NodesAnalyzed)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (jF *JsonFormatter) IsInterfaceNil() bool {
	return jF == nil
}
