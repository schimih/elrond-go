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
	mut               sync.Mutex
	discoveredTargets []DiscoveredTarget
	discoverer        Discoverer
	scannerFactory    ScannerFactory
	parserFactory     ParserFactory
	analysisType      utils.AnalysisType
	managerCommand    int
}

var log = logger.GetOrCreate("vat")

// NewAnalyzer creates a new analyzer used for discovery and parsing activities
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
	a.managerCommand = NoCommand
	a.discoveredTargets = make([]DiscoveredTarget, 0)
	a.scannerFactory = sf
	a.parserFactory = pf

	return a, nil
}

// StartJob discovers new targets and start the analysis job
func (a *Analyzer) StartJob(analysisType utils.AnalysisType) (scanResults []scan.ScannedTarget) {
	a.mut.Lock()
	defer a.mut.Unlock()

	a.discoverTargets()
	// get command from manager
	a.analysisType = analysisType
	nmapScanResults := a.deployAnalysisWorkers()
	p := a.parserFactory.CreateParser(nmapScanResults, analysisType)
	return p.Parse()
}

func (a *Analyzer) discoverTargets() {
	a.discoveredTargets = a.discoverer.DiscoverNewTargets(a.discoveredTargets)
}

func (a *Analyzer) deployAnalysisWorkers() (work [][]byte) {
	scanResults := make([][]byte, 0)
	var wg sync.WaitGroup
	for _, h := range a.discoveredTargets {
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

// this is concurrent safe because a target is not accessed by two concurrent workers
func (a *Analyzer) worker(h *DiscoveredTarget) (rawScanResults []byte) {
	s := a.scannerFactory.CreateScanner(h.Address, utils.AnalysisType(a.analysisType))

	log.Info("Starting scan for:", "address", h.Address)
	// Run the scan
	rawResult, err := s.Scan()
	if err != nil {
		log.Error("Scan failed because %e", err)
	}

	a.changeTargetStatus(h, utils.SCANNED)

	log.Info("Scanning done for target:", "address", h.Address)
	return rawResult
}

func (a *Analyzer) changeTargetStatus(h *DiscoveredTarget, status utils.TargetStatus) {
	h.Status = status
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Analyzer) IsInterfaceNil() bool {
	return a == nil
}
