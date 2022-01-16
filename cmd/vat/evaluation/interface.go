package evaluation

// type EvaluatorFactory interface {
// 	CreateEvaluator(node Node, evaluationType int) Evaluator
// 	IsInterfaceNil() bool
// }

type Evaluator interface {
	Evaluate() EvaluationTarget
	IsInterfaceNil() bool
}
