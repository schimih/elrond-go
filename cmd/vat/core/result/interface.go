package result

type ResultManager interface {
	Process()
	Evaluate()
	Export()
}

type AssessmentRunner interface {
	Assess()
	Name()
	IsInterfaceNil()
}
