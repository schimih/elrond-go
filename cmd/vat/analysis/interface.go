package analysis

import (
	"github.com/elrond-go/cmd/vat/core"
	"github.com/elrond-go/cmd/vat/scan"
)

type Discoverer interface {
	DiscoverNewTargets(existingTargets []DiscoveredTarget) (targets []DiscoveredTarget)
	IsInterfaceNil() bool
}

type ScannerFactory interface {
	CreateScanner(target string, analysisType core.AnalysisType) scan.Scanner
	IsInterfaceNil() bool
}
type ParserFactory interface {
	CreateParser(input [][]byte, grammar core.AnalysisType) scan.Parser
	IsInterfaceNil() bool
}
