package scan

import (
	"os/exec"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	core "github.com/elrond-go/cmd/vat/core"
)

var log = logger.GetOrCreate("vat")

type NmapScanner struct {
	mutScanner sync.Mutex
	Name       string
	Target     string
	Status     core.ScannerStatus
	Cmd        string
}

func (s *NmapScanner) preScan() {
	s.Status = core.IN_PROGRESS
}
func (s *NmapScanner) postScan() {
	s.Status = core.FINISHED
}

// Run nmap scan
func (s *NmapScanner) Scan() (res []byte, err error) {
	s.mutScanner.Lock()
	defer s.mutScanner.Unlock()
	s.preScan()
	res, err = shellCmd(s.Cmd)
	if err != nil {
		return nil, err
	}

	s.postScan()
	return res, nil
}

func shellCmd(cmd string) (res []byte, err error) {

	res, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *NmapScanner) IsInterfaceNil() bool {
	return s == nil
}
