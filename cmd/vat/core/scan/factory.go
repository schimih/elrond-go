package scan

import "github.com/elrond-go/cmd/vat/core/result"

type NmapScannerFactory struct {
}

func NewNmapScannerFactory() *NmapScannerFactory {
	return &NmapScannerFactory{}
}

func (factory *NmapScannerFactory) CreateScanner(target string, analysisType string) Scanner {
	switch analysisType {
	case "TCP-ELROND":
		return &nmapScanner{
			name:   "TCP-ELROND",
			target: target,
			status: result.NOT_STARTED,
			cmd:    constructCmd(target, NMAP_TCP_ELROND),
		}
	case "TCP-WEB":
		return &nmapScanner{
			name:   "TCP-WEB",
			target: target,
			status: result.NOT_STARTED,
			cmd:    constructCmd(target, NMAP_TCP_WEB),
		}
	case "TCP-SSH":
		return &nmapScanner{
			name:   "TCP-SSH",
			target: target,
			status: result.NOT_STARTED,
			cmd:    constructCmd(target, NMAP_TCP_SSH),
		}
	case "TCP-ALL":
		return &nmapScanner{
			name:   "TCP-TOP",
			target: target,
			status: result.NOT_STARTED,
			cmd:    constructCmd(target, NMAP_TCP_STANDARD),
		}
	default:
		return nil
	}
}

func (factory *NmapScannerFactory) CreateElrondScanner(target string) Scanner {
	return &nmapScanner{
		name:   "TCP-ELROND",
		target: target,
		status: result.NOT_STARTED,
		cmd:    constructCmd(target, NMAP_TCP_ELROND),
	}
}

func (factory *NmapScannerFactory) CreateWebScanner(target string) Scanner {
	return &nmapScanner{
		name:   "TCP-WEB",
		target: target,
		status: result.NOT_STARTED,
		cmd:    constructCmd(target, NMAP_TCP_WEB),
	}
}

func (factory *NmapScannerFactory) CreateSSHScanner(target string) Scanner {
	return &nmapScanner{
		name:   "TCP-SSH",
		target: target,
		status: result.NOT_STARTED,
		cmd:    constructCmd(target, NMAP_TCP_SSH),
	}
}

func (factory *NmapScannerFactory) CreateTop1000Scanner(target string) Scanner {
	return &nmapScanner{
		name:   "TCP-TOP",
		target: target,
		status: result.NOT_STARTED,
		cmd:    constructCmd(target, NMAP_TCP_STANDARD),
	}
}
