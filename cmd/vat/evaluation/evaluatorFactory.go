package evaluation

import (
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

// initial rating
const initialRating = 100

type EvaluatorFactory struct {
}

func NewEvaluatorFactory() EvaluatorFactory {
	return EvaluatorFactory{}
}

func (eF *EvaluatorFactory) CreateEvaluator(address string, ports []scan.Port, analysisType utils.AnalysisType) Evaluator {
	return &EvaluatedTarget{
		Address:       address,
		Ports:         ports,
		Status:        string(utils.NEW),
		Score:         initialRating,
		SecurityLevel: utils.HIGH,
		Judgements:    make([]string, 0),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (eF *EvaluatorFactory) IsInterfaceNil() bool {
	return eF == nil
}
