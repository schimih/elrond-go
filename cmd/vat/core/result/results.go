package result

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/jinzhu/gorm"
)

type Step int

const (
	NOT_DEFINED Step = iota
	IMPORTED         // targets
	NEW              // hosts
	SCANNED          // hosts
)

const (
	NULL = iota
	NOT_STARTED
	IN_PROGRESS
	FAILED
	DONE
	FINISHED
)

type ResultsContainer struct {
	Results []Peer
}

var log = logger.GetOrCreate("vat")

func NewResultsContainer(messenger p2p.Messenger) *ResultsContainer {
	r := &ResultsContainer{}
	r.Results = make([]Peer, len(messenger.ConnectedAddresses()))
	return r
}

func (p *Port) GetPeer(db *gorm.DB) *Peer {
	// peer := &Peer{}
	// db.Where("id = ?", p.PeerID).Find(&peer)
	return nil
}

func GetAllPeers(db *gorm.DB) []Peer {
	//var err error
	peers := []Peer{}
	// check if empty.
	db.Find(&peers)
	if len(peers) == 0 {
		log.Error("DB EMPTY. Run 'discover_peer' first or load a json with peers")
		return nil
	}
	return peers
}

func (h *Peer) GetPorts(db *gorm.DB) []Port {
	ports := []Port{}
	db.Where("peer_id = ?", h.ID).Find(&ports)
	return ports
}
