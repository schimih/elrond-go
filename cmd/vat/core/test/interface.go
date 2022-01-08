package test

type TestRunner interface {
	Run()
	GetTargetsList() []Target
}
