package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/display"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	factoryMarshalizer "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/epochStart/bootstrap/disabled"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/VAT/cmd/vat/core/result"
	"github.com/VAT/cmd/vat/core/test"
	"github.com/urfave/cli"
)

type cfg struct {
	vulTestType string
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
	// p2pSeed defines a flag to be used as a seed when generating P2P credentials. Useful for seed nodes.
	testType = cli.StringFlag{
		Name:        "test-type",
		Usage:       "P2P seed will be used when generating credentials for p2p component. Can be any string.",
		Value:       "full",
		Destination: &argsConfig.vulTestType,
	}
	argsConfig           = &cfg{}
	p2pConfigurationFile = "./config/p2p.toml"
)

var log = logger.GetOrCreate("vat")

func main() {

	log.Info("Starting VAT")
	app := cli.NewApp()
	cli.AppHelpTemplate = vulTemplate
	app.Name = "Vulnerability Analysis Tool"
	app.Version = "v0.0.1"
	app.Usage = "This tool will be used for security checks on Elrond EcoSystem (v0.0.1 - portscanner and ssh access)"
	app.Flags = []cli.Flag{
		testType,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		return startVulnerabilityAnalysis(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
func startVulnerabilityAnalysis(ctx *cli.Context) error {
	var err error
	var vulTest string

	// Start seednode
	messenger, err := startSeedNode()
	if err != nil {
		return err
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Analysis is now running")
	mainLoop(messenger, sigs, vulTest)

	return nil
}

func mainLoop(messenger p2p.Messenger, stop chan os.Signal, vulTest string) {
	t := test.NewTest(messenger, vulTest)
	r := result.NewResultsContainer(messenger)

	for {
		select {
		case <-stop:
			log.Info("terminating at user's signal...")
			return
		case <-time.After(time.Second * 5):
			t.UpdateTargetsList()
			r.Evaluate(r.Process(t.Run(), vulTest))
			//displayMessengerInfo(messenger)
			log.Info("Added targets", "targets", len(t.TargetsList))
		}
	}
}

func loadMainConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createNode(p2pConfig config.P2PConfig, marshalizer marshal.Marshalizer) (p2p.Messenger, error) {
	arg := libp2p.ArgsNetworkMessenger{
		Marshalizer:          marshalizer,
		ListenAddress:        libp2p.ListenAddrWithIp4AndTcp,
		P2pConfig:            p2pConfig,
		SyncTimer:            &libp2p.LocalSyncTimer{},
		PreferredPeersHolder: disabled.NewPreferredPeersHolder(),
		NodeOperationMode:    p2p.NormalOperation,
	}

	return libp2p.NewNetworkMessenger(arg)
}

func startSeedNode() (messenger p2p.Messenger, err error) {

	generalConfig, err := loadMainConfig("./config/config.toml")
	if err != nil {
		return nil, err
	}

	internalMarshalizer, err := factoryMarshalizer.NewMarshalizer(generalConfig.Marshalizer.Type)
	if err != nil {
		return nil, fmt.Errorf("error creating marshalizer (internal): %s", err.Error())
	}

	log.Info("Starting Seed Node")

	p2pConfig, err := common.LoadP2PConfig(p2pConfigurationFile)
	if err != nil {
		return nil, err
	}

	messenger, err = createNode(*p2pConfig, internalMarshalizer)
	if err != nil {
		return nil, err
	}

	log.Info("Starting Bootstrap")
	err = messenger.Bootstrap()
	if err != nil {
		return nil, err
	}

	return messenger, nil
}

func displayMessengerInfo(messenger p2p.Messenger) {
	headerSeedAddresses := []string{"Seednode addresses:"}
	addresses := make([]*display.LineData, 0)

	for _, address := range messenger.Addresses() {
		addresses = append(addresses, display.NewLineData(false, []string{address}))
	}

	tbl, _ := display.CreateTableString(headerSeedAddresses, addresses)
	log.Info("\n" + tbl)

	mesConnectedAddrs := messenger.ConnectedAddresses()
	sort.Slice(mesConnectedAddrs, func(i, j int) bool {
		return strings.Compare(mesConnectedAddrs[i], mesConnectedAddrs[j]) < 0
	})

	log.Info("known peers", "num peers", len(messenger.Peers()))
	headerConnectedAddresses := []string{fmt.Sprintf("Seednode is connected to %d peers:", len(mesConnectedAddrs))}
	connAddresses := make([]*display.LineData, len(mesConnectedAddrs))

	for idx, address := range mesConnectedAddrs {
		connAddresses[idx] = display.NewLineData(false, []string{address})
	}

	tbl2, _ := display.CreateTableString(headerConnectedAddresses, connAddresses)
	log.Info("\n" + tbl2)
}
