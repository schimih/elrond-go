package evaluation

import "github.com/elrond-go/cmd/vat/utils"

const (
	initialRating = 100.0
)

type EvaluatorFactory struct {
}

func NewEvaluatorFactory() EvaluatorFactory {
	return EvaluatorFactory{}
}

func (ef *EvaluatorFactory) CreateEvaluator(node Node, evaluationType utils.AnalysisType) Evaluator {
	return &EvaluationResult{
		Node:           node,
		EvaluationType: evaluationType,
		Score:          initialRating,
		Judgment:       make([]string, 0),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (ef *EvaluatorFactory) IsInterfaceNil() bool {
	return ef == nil
}
