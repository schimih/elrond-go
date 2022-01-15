package factory

import (
	"fmt"

	"github.com/elrond-go/cmd/vat/core/scan"
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

const (
	TCP_ELROND = iota
	TCP_OUTSIDE_ELROND
	TCP_WEB
	TCP_SSH
	TCP_FULL
	TCP_STANDARD
)

type NmapScannerFactory struct {
}

func NewNmapScannerFactory() *NmapScannerFactory {
	return &NmapScannerFactory{}
}

func (factory *NmapScannerFactory) CreateScanner(target string, analysisType int) scan.Scanner {
	switch analysisType {
	case TCP_ELROND:
		return &scan.ArgNmapScanner{
			Name:   "TCP-ELROND",
			Target: target,
			Status: scan.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_ELROND),
		}
	case TCP_WEB:
		return &scan.ArgNmapScanner{
			Name:   "TCP-WEB",
			Target: target,
			Status: scan.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_WEB),
		}
	case TCP_SSH:
		return &scan.ArgNmapScanner{
			Name:   "TCP-SSH",
			Target: target,
			Status: scan.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_SSH),
		}
	case TCP_STANDARD:
		return &scan.ArgNmapScanner{
			Name:   "TCP-TOP",
			Target: target,
			Status: scan.NOT_STARTED,
			Cmd:    constructCmd(target, NMAP_TCP_STANDARD),
		}
	default:
		return nil
	}
}

func constructCmd(target string, args string) string {
	return fmt.Sprintf("nmap %s %s -oX -", args, target)
}

// IsInterfaceNil returns true if there is no value under the interface
func (factory *NmapScannerFactory) IsInterfaceNil() bool {
	return factory == nil
}
