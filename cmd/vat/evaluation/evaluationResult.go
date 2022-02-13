package evaluation

import "github.com/elrond-go/cmd/vat/core"

// initial rating
const initialRating = 100

type evaluationResult struct {
	score          int
	securityLevel  core.SecureLevel
	judgements     []string
	evaluationType core.EvaluationType
}

func newEvaluationResult(evaluationType core.EvaluationType) (evaluation evaluationResult) {
	return evaluationResult{
		score:          initialRating,
		securityLevel:  core.HIGH,
		judgements:     make([]string, 0),
		evaluationType: evaluationType,
	}
}
