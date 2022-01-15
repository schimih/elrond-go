package analysis

import "github.com/elrond-go/cmd/vat/core/scan"

type Discoverer interface {
	DiscoverNewTargets(existingTargets []Target) (targets []Target)
	IsInterfaceNil() bool
}

type ScannerFactory interface {
	CreateScanner(target string, analysisType int) scan.Scanner
	IsInterfaceNil() bool
}
type ParserFactory interface {
	CreateParser(input [][]byte, grammar int) scan.Parser
	IsInterfaceNil() bool
}
