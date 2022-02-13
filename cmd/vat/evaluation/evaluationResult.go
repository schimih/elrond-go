package evaluation

import "github.com/elrond-go/cmd/vat/utils"

type evaluationResult struct {
	score          int
	securityLevel  utils.SecureLevel
	judgements     []string
	evaluationType utils.EvaluationType
}

func newEvaluationResult(evaluationType utils.EvaluationType) (evaluation evaluationResult) {
	return evaluationResult{
		score:          initialRating,
		securityLevel:  utils.HIGH,
		judgements:     make([]string, 0),
		evaluationType: evaluationType,
	}
}
