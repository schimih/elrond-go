package scan

import (
	"fmt"
	"os/exec"
	"strings"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elrond-go/cmd/vat/core/result"
	go_nmap "github.com/lair-framework/go-nmap"
)

/*
-Pn --skip the ping test and simply scan every target host provided.
-sS --stealth scan,fastest way to scan ports of the most popular protocol (TCP).
-pn --port to be scanned.
-sC --
*/
var NMAP_TCP_ELROND = "-Pn -sS -p37373-38383"
var NMAP_TCP_OUTSIDE_ELROND = "-Pn -sS -p-37372,38384-"
var NMAP_TCP_WEB = "-Pn -p80,8080,280,443" // added: http-mgmt (280), https (443)
var NMAP_TCP_SSH = "-Pn -p22"
var NMAP_TCP_FULL = "-Pn -sS -A -p-"
var NMAP_TCP_STANDARD = "--randomize-hosts -Pn -sS -A -T4 -g53 --top-ports 1000"
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

// Constructor for NmapScan
func NewNmapScanner(name string, target string, nmapArgs string) *nmapScanner {
	return &nmapScanner{
		name:   name,
		target: target,
		status: result.NOT_STARTED,
		cmd:    constructCmd(target, nmapArgs),
	}
}

func (s *nmapScanner) preScan() {
	s.status = result.IN_PROGRESS
}
func (s *nmapScanner) postScan() {
	s.status = result.FINISHED
}

func constructCmd(target string, args string) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// Run nmap scan
func (s *nmapScanner) Scan() (res *go_nmap.NmapRun) {
	s.preScan()
	// Run nmap
	res, err := shellCmd(s.cmd)
	if err != nil {
		s.status = result.FAILED
	}

	s.postScan()
	return res
}

func shellCmd(cmd string) (result *go_nmap.NmapRun, err error) {
	// Prepare Nmap
	res, err := exec.Command("sh", "-c", cmd).Output()
	// Run Nmap
	if err != nil {
		return nil, err
	}
	result, _ = go_nmap.Parse(res)
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			log.Error(err.Error())
		}
		return result, err
	}
	return result, err
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *nmapScanner) IsInterfaceNil() bool {
	return s == nil
}
