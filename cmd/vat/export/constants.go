package export

type OutputType int

const (
	Table    OutputType = 0
	JSON     OutputType = 1
	XML      OutputType = 2
	GIN      OutputType = 3
	JustScan OutputType = 4
)

const JsonFilePath = "./cmd/vat/AnalysisResult.json"
const XMLFilePath = "./cmd/vat/AnalysisResult.xml"
