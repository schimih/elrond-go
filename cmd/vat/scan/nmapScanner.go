package scan

import (
	"os/exec"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/utils"
)

var log = logger.GetOrCreate("vat")

type ArgNmapScanner struct {
	Name   string
	Target string
	Status utils.ScannerStatus
	Cmd    string
}

type nmapScanner struct {
	name   string
	target string
	status utils.ScannerStatus
	cmd    string
}

func (s *ArgNmapScanner) preScan() {
	s.Status = utils.IN_PROGRESS
}
func (s *ArgNmapScanner) postScan() {
	s.Status = utils.FINISHED
}

// Run nmap scan
func (s *ArgNmapScanner) Scan() (res []byte, err error) {
	s.preScan()
	res, error := shellCmd(s.Cmd)
	if error != nil {
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
func (s *nmapScanner) IsInterfaceNil() bool {
	return s == nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *ArgNmapScanner) IsInterfaceNil() bool {
	return s == nil
}
