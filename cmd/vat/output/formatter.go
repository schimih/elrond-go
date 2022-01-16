package output

import (
	"github.com/elrond-go/cmd/vat/evaluation"
	"github.com/elrond-go/cmd/vat/utils"
)

func CreateFormatter(formatType utils.OutputType, evaluationReport []evaluation.EvaluationTarget) Formatter {
	switch formatType {
	case utils.Table:
		return &TableFormatter{
			EvaluationReport: evaluationReport,
		}
	case utils.JSON:
		return &JsonFormatter{
			EvaluationReport: evaluationReport,
		}
	case utils.XML:
		return &JsonFormatter{
			EvaluationReport: evaluationReport,
		}
	default:
		return nil
	}
}
