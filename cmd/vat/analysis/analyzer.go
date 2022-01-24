package analysis

import (
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type Analyzer struct {
	DiscoveredTargets []DiscoveredTarget
	discoverer        Discoverer
	scannerFactory    ScannerFactory
	parserFactory     ParserFactory
	AnalysisType      utils.AnalysisType
	ManagerCommand    int
}

var log = logger.GetOrCreate("vat")

func NewAnalyzer(discoverer Discoverer, sf ScannerFactory, pf ParserFactory) (*Analyzer, error) {
	if check.IfNil(discoverer) {
		return nil, fmt.Errorf("Discoverer needed")
	}

	if check.IfNil(sf) {
		return nil, fmt.Errorf("ScannerFactory needed")
	}

	if check.IfNil(pf) {
		return nil, fmt.Errorf("ParserFactory needed")
	}

	a := &Analyzer{}
	a.discoverer = discoverer
	a.ManagerCommand = NO_COMMAND
	a.DiscoveredTargets = make([]DiscoveredTarget, 0)
	a.scannerFactory = sf
	a.parserFactory = pf

	return a, nil
}

func (a *Analyzer) DiscoverTargets() {
	a.DiscoveredTargets = a.discoverer.DiscoverNewTargets(a.DiscoveredTargets)
}

func (a *Analyzer) AnalyzeNewlyDiscoveredTargets(analysisType utils.AnalysisType) (scanResults []scan.ScannedTarget) {
	// get command from manager
	a.AnalysisType = analysisType
	nmapScanResults := a.deployAnalysisWorkers()
	p := a.parserFactory.CreateParser(nmapScanResults, analysisType)
	return p.Parse()
}

func (a *Analyzer) deployAnalysisWorkers() (work [][]byte) {
	scanResults := make([][]byte, 0)
	var wg sync.WaitGroup
	for _, h := range a.DiscoveredTargets {
		if (h.ActualStatus() == New) || (h.ActualStatus() == Expired) {
			wg.Add(1)
			temp := h
			go func() {
				defer wg.Done()
				scanResults = append(scanResults, a.worker(&temp))
			}()
		}
	}
	wg.Wait()
	return scanResults
}

func (a *Analyzer) worker(h *DiscoveredTarget) (rawScanResults []byte) {
	s := a.scannerFactory.CreateScanner(h.Address, utils.AnalysisType(a.AnalysisType))

	log.Info("Starting scan for:", "address", h.Address)
	// Run the scan
	rawResult, err := s.Scan()
	if err != nil {
		log.Error("Scan failed because %e", err)
	}

	log.Info("Scanning done for target:", "address", a.changeTargetStatus(h.Address, utils.SCANNED))
	return rawResult
}

func (a *Analyzer) changeTargetStatus(address string, status utils.TargetStatus) string {
	for idx, _ := range a.DiscoveredTargets {
		if address == a.DiscoveredTargets[idx].Address {
			a.DiscoveredTargets[idx].Status = status
			break
		}
	}
	return address
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Analyzer) IsInterfaceNil() bool {
	return a == nil
}
