package analysis

import (
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	go_nmap "github.com/lair-framework/go-nmap"
)

type Analyzer struct {
	Targets      []Target
	discoverer   Discoverer
	scanner      ScannerFactory
	AnalysisType string
}

// TargetStatus represents a target's state.
type TargetStatus string

// Enumerates the different possible state values.
const (
	New     TargetStatus = "NEW"
	Scanned TargetStatus = "SCANNED"
)

// Status returns the status of a target.
func (t Target) ActualStatus() TargetStatus {
	return TargetStatus(t.Status)
}

var log = logger.GetOrCreate("vat")

func NewAnalyzer(discoverer Discoverer, sf ScannerFactory, AnalysisType string) (*Analyzer, error) {
	if check.IfNil(discoverer) {
		return nil, fmt.Errorf("Discoverer needed")
	}

	a := &Analyzer{}
	a.discoverer = discoverer
	a.AnalysisType = AnalysisType
	a.Targets = make([]Target, 0)
	a.scanner = sf

	return a, nil
}

func (a *Analyzer) DiscoverNewPeers() {
	a.Targets = a.discoverer.DiscoverNewTargets(a.Targets)
}

func (a *Analyzer) StartScan() (testresults []*go_nmap.NmapRun) {
	nmapScanResults := make([]*go_nmap.NmapRun, 0)
	var wg sync.WaitGroup
	for _, h := range a.Targets {
		if h.ActualStatus() == "NEW" {
			wg.Add(1)
			temp := h
			go func() {
				defer wg.Done()
				nmapScanResults = append(nmapScanResults, a.worker(&temp))
			}()
			//wg.Wait() //=
		} else {
			log.Info("Target already scanned", "address", h.Address)
		}
	}
	wg.Wait()
	return nmapScanResults
}

func (a *Analyzer) worker(h *Target) (rawResult *go_nmap.NmapRun) {

	s := a.scanner.CreateScanner(h.Address, a.AnalysisType)
	log.Info("Starting scan for:", "address", h.Address)
	// Run the scan
	res := s.Scan()
	if res != nil {
		a.ProcessScanResult(h.Address)
	}
	log.Info("Scanning done for target:", "address", h.Address)
	return res
}

func (a *Analyzer) ProcessScanResult(address string) error {
	// Each time a test worker finishes test, change status to Scanned in targets list
	for idx, target := range a.Targets {
		if target.Address == address {
			a.Targets[idx].Status = "SCANNED"
			log.Info("changed state to SCANNED for ", "address", target.Address)
		}
	}
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Analyzer) IsInterfaceNil() bool {
	return a == nil
}
