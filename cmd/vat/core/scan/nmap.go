package scan

import (
	"os/exec"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

const (
	NULL = iota
	NOT_STARTED
	IN_PROGRESS
	FAILED
	DONE
	FINISHED
)

var log = logger.GetOrCreate("vat")

type ArgNmapScanner struct {
	Name   string
	Target string
	Status int
	Cmd    string
}

type nmapScanner struct {
	name   string
	target string
	status int
	cmd    string
}

func (s *ArgNmapScanner) preScan() {
	s.Status = IN_PROGRESS
}
func (s *ArgNmapScanner) postScan() {
	s.Status = FINISHED
}

// Run nmap scan
func (s *ArgNmapScanner) Scan() (res []byte) {
	s.preScan()
	// Run nmap
	res, err := shellCmd(s.Cmd)
	if err != nil {
		s.Status = FAILED
	}

	s.postScan()
	return res
}

func shellCmd(cmd string) (res []byte, err error) {
	// Prepare Nmap
	res, err = exec.Command("sh", "-c", cmd).Output()
	// Run Nmap
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
