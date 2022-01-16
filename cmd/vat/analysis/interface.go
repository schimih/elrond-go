package analysis

import (
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type Discoverer interface {
	DiscoverNewTargets(existingTargets []Target) (targets []Target)
	IsInterfaceNil() bool
}

type ScannerFactory interface {
	CreateScanner(target string, analysisType utils.AnalysisType) scan.Scanner
	IsInterfaceNil() bool
}
type ParserFactory interface {
	CreateParser(input [][]byte, grammar utils.AnalysisType) scan.Parser
	IsInterfaceNil() bool
}
