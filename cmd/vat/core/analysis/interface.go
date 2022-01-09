package analysis

import "github.com/elrond-go/cmd/vat/core/scan"

type AnalyzeRunner interface {
	Run()
	IsInterfaceNil() bool
}

type Discoverer interface {
	DiscoverNewTargets(existingTargets []Target) (targets []Target)
	IsInterfaceNil() bool
}

type ScannerFactory interface {
	CreateScanner(target string, analysisType string) scan.Scanner
}
