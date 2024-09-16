package evaluation

type Evaluator interface {
	Evaluate() EvaluatedTarget
	IsInterfaceNil() bool
}
