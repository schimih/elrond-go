package manager

type Formatter interface {
	Output(rankedReport RankedReport)
	IsInterfaceNil() bool
}
