package scan

type Scanner interface {
	Scan() (res []byte, err error)
	IsInterfaceNil() bool
}

type Parser interface {
	Parse() (parsingResults []ScannedTarget)
	IsInterfaceNil() bool
}

type Dialer interface {
	Dial() error
	IsInterfaceNil() bool
}
