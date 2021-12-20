package main

import (
	"fmt"
	"os"
	"strconv"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/VAT/cmd/vul/core/scan"
	"github.com/VAT/cmd/vul/core/utils"
	"github.com/urfave/cli"
)

type cfg struct {
	peersScanNumber   int
	peersScanAll      bool
	peersFileLoad     string
	peersFileShow     string
	portScanElrond    bool
	portScan8080      bool
	portScanSSH       bool
	portScanAll       bool
	portScanLoad      string
	ssh22             bool
	sshUsrPsw         bool
	utilsOutputFolder string
	showSSHUsrPsw     string
}

const (
	defaultLogsPath     = "logs"
	logFilePrefix       = "elrond-seed"
	filePathPlaceholder = "[path]"
)

var (
	vulTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}
 VERSION:
	{{.Version}}
	{{end}}
 `
	peersScanNumber = cli.IntFlag{
		Name:        "number",
		Usage:       "Number of unique peers to discover",
		Value:       30,
		Destination: &argsConfig.peersScanNumber,
	}
	peersScanAll = cli.BoolFlag{
		Name:  "all",
		Usage: "Discover all possible peers",
	}

	peersFileLoad = cli.StringFlag{
		Name:        "load",
		Usage:       "Introduce a path to an already populated json file for vul analysis",
		Destination: &argsConfig.peersFileLoad,
	}

	peersFileShow = cli.StringFlag{
		Name:  "show",
		Usage: "Log peers data from an already populated json file",
	}

	portScanElrond = cli.BoolFlag{
		Name:  "tcp-elrond",
		Usage: "Verifying if just the elrond port list (37373-38383) is open",
	}

	portScanWEB = cli.BoolFlag{
		Name:  "tcp-web",
		Usage: "Verifying if port 80 or 8080 is accessible",
	}

	portScanSSH = cli.BoolFlag{
		Name:  "tcp-ssh",
		Usage: "Verifying if the ssh port 22 is accessible",
	}

	portScanALL = cli.BoolFlag{
		Name:  "all",
		Usage: "Checks everything",
	}

	portScanLoad = cli.StringFlag{
		Name:        "load",
		Usage:       "Introduce a path to an already populated json file for results analysis or update",
		Destination: &argsConfig.portScanLoad,
	}

	ssh22 = cli.StringFlag{
		Name:  "22",
		Usage: "Check only for ssh - if the ssh port:22 is accessible",
	}

	sshUsrPsw = cli.StringFlag{
		Name:  "usr_psw",
		Usage: "Check if username or password could be sent to peer",
	}

	utilsOutputFolder = cli.StringFlag{
		Name:  "output",
		Usage: "Set Output folder",
	}

	peersEval = cli.BoolFlag{
		Name:  "start_eval",
		Usage: "Start evaluation using data already populated in DB",
	}

	peersScanShow = cli.BoolFlag{
		Name:  "show_peers",
		Usage: "Show scan result from DB",
	}

	portScanShow = cli.BoolFlag{
		Name:  "show_ports",
		Usage: "Show scan result from DB",
	}

	argsConfig           = &cfg{}
	p2pConfigurationFile = "./config/p2p.toml"
)

var log = logger.GetOrCreate("vat")

func main() {
	//init
	log.Info("Starting VAT")
	utils.InitConfig()
	app := cli.NewApp()
	cli.AppHelpTemplate = vulTemplate
	app.Name = "Vulnerability Analysis Tool"
	app.Version = "v0.0.1"
	app.Usage = "This tool will be used for security checks on Elrond EcoSystem (v0.0.1 - portscanner and ssh access)"
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}
	//maybe args should have been used.
	app.Commands = []cli.Command{
		{
			Name:  "discover_peer",
			Usage: "Start Peer Discovery and save results to DB and/or .json, default 30",
			Flags: []cli.Flag{
				peersScanNumber, // done - without refactor.
				peersScanAll,    // todo:
				peersFileLoad,   // done
				peersFileShow,   // todo
			},
			Action: func(c *cli.Context) error {
				ScanForPeers(c)
				return nil
			},
		},
		{
			Name:  "scan_port",
			Usage: "Perform different port scans and update DB and/or .json ",
			Flags: []cli.Flag{
				portScanElrond, // done
				portScanWEB,    // done
				portScanSSH,    // done
				portScanALL,    // done
				portScanLoad,   // todo
			},
			Action: func(c *cli.Context) error {
				ScanPort(c)
				return nil
			},
		},
		{
			Name:  "scan_ssh",
			Usage: "SSH Service",
			Flags: []cli.Flag{
				ssh22,     // done, make nmap commd
				sshUsrPsw, // how? with a txt as input for forcing
			},
			Action: func(c *cli.Context) error {
				ScanSSH(c)
				return nil
			},
		},
		{
			Name:  "utils",
			Usage: "General configuration for analysis",
			Flags: []cli.Flag{
				peersEval,     // todo
				portScanShow,  // done
				peersScanShow, // todo
			},
			Action: func(c *cli.Context) error {
				Utils(c)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func ScanForPeers(ctx *cli.Context) {
	//check what to do
	if ctx.IsSet(peersScanNumber.Name) {
		reqPeers := uint(ctx.Int(peersScanNumber.Name))

		log.Info("Starting discovery for", "peers", reqPeers)
		scan.ScanForPeers(reqPeers)
	}
	if ctx.IsSet(peersFileLoad.Name) {
		path := fmt.Sprintf(ctx.String(peersFileLoad.Name))

		log.Info("Loading json file:", "path", path)
		utils.LoadFile(path)
	}
	//more to add
}

func ScanPort(ctx *cli.Context) {
	// ugly,
	// refactor - use enums
	var ScanType string
	if ctx.IsSet(portScanElrond.Name) {
		ScanType = "TCP-ELROND"
	}
	if ctx.IsSet(portScanWEB.Name) {
		ScanType = "TCP-WEB"
	}
	if ctx.IsSet(portScanSSH.Name) {
		ScanType = "TCP-SSH"
	}
	if ctx.IsSet(portScanALL.Name) {
		ScanType = "TCP-ALL"
	}
	if ctx.IsSet(portScanLoad.Name) {
		ScanType = "TCP-LOAD"
	}
	scan.ScanPort(ScanType)
	log.Info("Scanning complete. In order to show results run 'scan_port -show")
	//more to add
}

// todo
func ScanSSH(ctx *cli.Context) {
	//check what to do
	if ctx.IsSet(peersScanNumber.Name) {
		requestedNumberPeers := ctx.GlobalString(peersScanNumber.Name)
		requestedPeers, _ := strconv.ParseUint(requestedNumberPeers, 10, 64)
		fmt.Println(requestedPeers)
		//scan.ScanForPeers(requestedPeers)
	}
}

// tofinish
func Utils(ctx *cli.Context) {
	//check what to do
	if ctx.IsSet(portScanShow.Name) {
		utils.ScanShow("ports")
	}
}
