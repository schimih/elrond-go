package scan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	factoryMarshalizer "github.com/ElrondNetwork/elrond-go-core/marshal/factory"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/epochStart/bootstrap/disabled"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/VAT/cmd/vul/core/model"
	"github.com/VAT/cmd/vul/core/utils"
)

var (
	p2pConfigurationFile = "./config/p2p.toml"
)

// Peer struct which contains an address
// a type, port number status.
type Peers struct {
	Peers []Peer
}
type Peer struct {
	ID             uint   `json:"id"`
	Protocol       string `json:"protocol"`
	Address        string `json:"address"`
	Standard       string `json:"standard"`
	ConnectionPort string `json:"connection-port"`
	Architecture   string `json:"architecture"`
	Messenger      string `json:"messenger-id"`
	Rating         int    `json:"rating"`
}

var log = logger.GetOrCreate("vat")

func ScanForPeers(requestedNumberPeers uint) error {

	var err error
	generalConfig, err := loadMainConfig("./config/config.toml")
	if err != nil {
		return err
	}

	internalMarshalizer, err := factoryMarshalizer.NewMarshalizer(generalConfig.Marshalizer.Type)
	if err != nil {
		return fmt.Errorf("error creating marshalizer (internal): %s", err.Error())
	}

	log.Info("Starting Seed Node")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	p2pConfig, err := common.LoadP2PConfig(p2pConfigurationFile)
	if err != nil {
		return err
	}

	messenger, err := createNode(*p2pConfig, internalMarshalizer)
	if err != nil {
		return err
	}

	log.Info("Starting Bootstrap")
	err = messenger.Bootstrap()
	if err != nil {
		return err
	}

	log.Info("Starting Peer Discovery")
	peersDiscoveryLoop(messenger, sigs, requestedNumberPeers)

	return nil
}

func peersDiscoveryLoop(mess p2p.Messenger, stop chan os.Signal, reqPeers uint) {
	var newUnqPeers uint
	var uniqueMasterSlice []Peer

	_ = peersDiscovery(mess)
	for newUnqPeers < reqPeers {
		time.Sleep(time.Second * 5)
		discoveredPeers := peersDiscovery(mess)
		uniqueMasterSlice = deduplicateAndStore(discoveredPeers, uniqueMasterSlice)
		newUnqPeers = uint(len(uniqueMasterSlice))
		log.Info("Discovered in total", "unique peers", newUnqPeers)
	}
	savePeersToJson(uniqueMasterSlice)
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

func peersDiscovery(mess p2p.Messenger) (foundPeersSlice []Peer) {

	mesConnectedAddrs := mess.ConnectedAddresses()
	// add host to db if new and return also the
	PeersSlice := createMasterSlice(mesConnectedAddrs)

	return PeersSlice
}

func createMasterSlice(connectedAddrs []string) (s []Peer) {
	// 143-150 refactor
	masterSlice := make([]Peer, len(connectedAddrs))
	for idx, address := range connectedAddrs {
		peerAddress := strings.Split(address, "/")
		masterSlice[idx].ID = uint(idx)
		masterSlice[idx].Protocol = peerAddress[1]
		masterSlice[idx].Address = peerAddress[2]
		masterSlice[idx].Standard = peerAddress[3]
		masterSlice[idx].ConnectionPort = peerAddress[4]
		masterSlice[idx].Architecture = peerAddress[5]
		masterSlice[idx].Messenger = peerAddress[6]
		masterSlice[idx].Rating = 100
	}

	return masterSlice
}

func deduplicateAndStore(discoveredSlice []Peer, uniqueSlice []Peer) (masterUniqueSlice []Peer) {
	//populate uniquepeers by looking into allPeers and copy the later to json
	//to refactor and use map
	masterUniqueSlice = uniqueSlice
	for _, peer := range discoveredSlice {
		exists := false
		for j, _ := range uniqueSlice {
			exists = reflect.DeepEqual(peer.Address, uniqueSlice[j].Address)
			if exists {
				break
			}
		}
		if !exists {
			//if unique peer, add to master slice + db
			masterUniqueSlice = append(masterUniqueSlice, peer)
			model.AddPeer(utils.Config.DB, peer.ID, peer.Protocol,
				peer.Standard, peer.Architecture, peer.Messenger,
				peer.Address, peer.Rating, model.NEW.String())
			log.Info("Added peer with", "address", peer.Address)
		}
	}
	return masterUniqueSlice
}

func savePeersToJson(InputUniqueSlice []Peer) {
	//marshal data for json
	//check the 2 write actions
	jsonData, _ := json.MarshalIndent(InputUniqueSlice, "", " ")
	_ = ioutil.WriteFile("peers.json", jsonData, 0644)

	path := "./cmd/vul/peers.json"
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return
	}
	log.Info("Peers list added to ", "path", path)
	// write to file, f.Write()
	f.Write(jsonData)
}
