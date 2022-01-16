package scan

type Scanner interface {
	Scan() []byte
	IsInterfaceNil() bool
}

type Parser interface {
	Parse() (parsingResults []AnalyzedTarget)
	IsInterfaceNil() bool
}
