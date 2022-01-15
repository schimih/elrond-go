package analysis

import (
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/core/scan"
)

type Analyzer struct {
	Targets        []Target
	discoverer     Discoverer
	scanner        ScannerFactory
	parser         ParserFactory
	AnalysisType   int
	ManagerCommand int
}

var log = logger.GetOrCreate("vat")

func NewAnalyzer(discoverer Discoverer, sf ScannerFactory, pf ParserFactory, analysisType int) (*Analyzer, error) {
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
	a.AnalysisType = analysisType
	a.ManagerCommand = NO_COMMAND
	a.Targets = make([]Target, 0)
	a.scanner = sf
	a.parser = pf

	return a, nil
}

func (a *Analyzer) DiscoverNewPeers() {
	a.Targets = a.discoverer.DiscoverNewTargets(a.Targets)
}

func (a *Analyzer) Analyze() (scanResults []scan.Peer) {
	nmapScanResults := a.deployAnalysisWorkers()
	p := a.parser.CreateParser(nmapScanResults, a.AnalysisType)
	return p.Parse()
}

func (a *Analyzer) deployAnalysisWorkers() (work [][]byte) {
	nmapScanResults := make([][]byte, 0)
	var wg sync.WaitGroup
	for _, h := range a.Targets {
		if (h.ActualStatus() == New) || (h.ActualStatus() == Expired) {
			wg.Add(1)
			temp := h
			go func() {
				defer wg.Done()
				nmapScanResults = append(nmapScanResults, a.worker(&temp))
			}()
		}
	}
	wg.Wait()
	return nmapScanResults
}

func (a *Analyzer) worker(h *Target) (scanRawResult []byte) {
	s := a.scanner.CreateScanner(h.Address, a.AnalysisType)

	log.Info("Starting scan for:", "address", h.Address)
	// Run the scan
	rawResult := s.Scan()

	log.Info("Scanning done for target:", "address", a.changeTargetStatus(h.Address, scan.SCANNED))
	return rawResult
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Analyzer) IsInterfaceNil() bool {
	return a == nil
}

func (a *Analyzer) changeTargetStatus(address string, status scan.TargetStatus) string {
	for idx, _ := range a.Targets {
		if address == a.Targets[idx].Address {
			a.Targets[idx].Status = status
			break
		}
	}
	return address
}
