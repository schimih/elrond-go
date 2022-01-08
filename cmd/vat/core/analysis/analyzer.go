package analysis

import (
	"fmt"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/core/scan"
	go_nmap "github.com/lair-framework/go-nmap"
)

type Analyzer struct {
	Targets      []Target
	discoverer   Discoverer
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

func NewAnalyzer(discoverer Discoverer, testType string) (*Analyzer, error) {
	if check.IfNil(discoverer) {
		return nil, fmt.Errorf("Discoverer needed")
	}
	a := &Analyzer{}
	a.discoverer = discoverer
	a.AnalysisType = testType
	a.Targets = make([]Target, 0)

	return a, nil
}

func (a *Analyzer) DiscoverNewPeers() {
	a.Targets = a.discoverer.DiscoverNewTargets(a.Targets)
}

func (t *Analyzer) Run() (NmapScanResults []*go_nmap.NmapRun) {
	NmapTestResults := make([]*go_nmap.NmapRun, 0)
	//force full test
	t.AnalysisType = "TCP-ALL"
	switch t.AnalysisType {
	case "TCP-ELROND":
		log.Info("Verifying if just the elrond port list (37373-38383) is open")
		NmapScanResults = t.StartScan(scan.NMAP_TCP_ELROND, NmapTestResults)
	case "TCP-WEB":
		log.Info("Verifying if port 80 or 8080 is accessible")
		NmapScanResults = t.StartScan(scan.NMAP_TCP_WEB, NmapTestResults)
	case "TCP-SSH":
		log.Info("Verifying if the ssh port 22 is accessible")
		NmapScanResults = t.StartScan(scan.NMAP_TCP_SSH, NmapTestResults)
	case "TCP-ALL":
		log.Info("Scans everything")
		NmapScanResults = t.StartScan(scan.NMAP_TCP_STANDARD, NmapTestResults)
	default:
		log.Error("Command unkown, no test will start", "command", t.AnalysisType)
		return
	}
	return NmapScanResults
}

func (t *Analyzer) StartScan(nmapArgs string, NmapTestResults []*go_nmap.NmapRun) (testresults []*go_nmap.NmapRun) {
	var wg sync.WaitGroup //=
	for _, h := range t.Targets {
		if h.ActualStatus() == "NEW" {
			wg.Add(1) //=

			go func() { //=
				defer wg.Done() //=
				temp := h
				NmapTestResults = append(NmapTestResults, t.worker(t.AnalysisType, &temp, nmapArgs))
			}()
			wg.Wait() //=
		} else {
			log.Info("Target already scanned", "address", h.Address)
		}
	}
	//wg.Wait() //=
	return NmapTestResults
}

func (t *Analyzer) worker(name string, h *Target, nmapArgs string) (rawResult *go_nmap.NmapRun) {

	s := scan.CreateNmapScanner(name, h.Address, nmapArgs)
	log.Info("Starting scan for:", "address", h.Address)
	// Run the scan
	res := s.RunNmap()
	if res != nil {
		// hosts - because nmap
		t.ProcessScanResult(h.Address)
	}
	log.Info("Scanning done for peer:", "address", h.Address)
	return res
}

func (t *Analyzer) ProcessScanResult(address string) error {
	// Each time a test worker finishes test, change status to Scanned in targets list
	for idx, target := range t.Targets {
		if target.Address == address {
			t.Targets[idx].Status = "SCANNED"
			log.Info("changed state to SCANNED for ", "address", target.Address)
		}
	}
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *Analyzer) IsInterfaceNil() bool {
	return a == nil
}
