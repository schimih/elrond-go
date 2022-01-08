package test

import (
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/VAT/cmd/vul/core/scan"
	go_nmap "github.com/lair-framework/go-nmap"
)

type Test struct {
	TargetsList []Target
	Messenger   p2p.Messenger
	TestType    string
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

func NewTest(messenger p2p.Messenger, testType string) *Test {
	t := &Test{}
	t.Messenger = messenger
	t.TestType = testType
	t.TargetsList = make([]Target, len(messenger.ConnectedAddresses()))
	return t
}

func (t *Test) Run() (NmapScanResults []*go_nmap.NmapRun) {
	NmapTestResults := make([]*go_nmap.NmapRun, 0)
	//force full test
	t.TestType = "TCP-ALL"
	switch t.TestType {
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
		log.Error("Command unkown, no test will start", "command", t.TestType)
		return
	}
	return NmapScanResults
}

func (t *Test) StartScan(nmapArgs string, NmapTestResults []*go_nmap.NmapRun) (testresults []*go_nmap.NmapRun) {
	var wg sync.WaitGroup //=
	for _, h := range t.TargetsList {
		if h.ActualStatus() == "NEW" {
			wg.Add(1) //=

			go func() { //=
				defer wg.Done() //=
				temp := h
				NmapTestResults = append(NmapTestResults, t.worker(t.TestType, &temp, nmapArgs))
			}()
			wg.Wait() //=
		} else {
			log.Info("Target already scanned", "address", h.Address)
		}
	}
	//wg.Wait() //=
	return NmapTestResults
}

func (t *Test) worker(name string, h *Target, nmapArgs string) (rawResult *go_nmap.NmapRun) {

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

func (t *Test) ProcessScanResult(address string) error {
	// Each time a test worker finishes test, change status to Scanned in targets list
	for idx, target := range t.TargetsList {
		if target.Address == address {
			t.TargetsList[idx].Status = "SCANNED"
			log.Info("changed state to SCANNED for ", "address", target.Address)
		}
	}
	return nil
}
