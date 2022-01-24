package manager

type Formatter interface {
	Output(rankedReport RankedReport) error
	IsInterfaceNil() bool
}
