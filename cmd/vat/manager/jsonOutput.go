package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/elrond-go/cmd/vat/utils"
)

type JsonFormatter struct {
}

func (jF *JsonFormatter) Output(rankedReport RankedReport) error {

	jsonData, _ := json.MarshalIndent(rankedReport, "", " ")
	_ = ioutil.WriteFile("AnalysisResult.json", jsonData, 0644)

	path := utils.JsonFilePath
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return fmt.Errorf("Could not open File")
	}
	log.Info("Peers list added to ", "path", path)
	f.Write(jsonData)
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (jF *JsonFormatter) IsInterfaceNil() bool {
	return jF == nil
}
