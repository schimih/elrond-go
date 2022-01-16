package output

type Formatter interface {
	GetOutput()
	IsInterfaceNil() bool
}
