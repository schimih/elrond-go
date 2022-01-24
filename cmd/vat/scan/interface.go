package scan

type Scanner interface {
	Scan() (res []byte, err error)
	IsInterfaceNil() bool
}

type Parser interface {
	Parse() (parsingResults []ScannedTarget)
	IsInterfaceNil() bool
}
