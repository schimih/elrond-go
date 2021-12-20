package model

import (
	"fmt"
	"os"
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Mutex sync.Mutex

var (
	lock sync.Mutex
)

type Step int

const (
	NOT_DEFINED Step = iota
	IMPORTED         // targets
	SWEEPED          // targets
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

/*==========
Model
todo: add description
=============*/
type Peer struct {
	ID           uint `gorm:"primary_key"`
	Protocol     string
	Address      string `gorm:"unique_index:idx_hostname_ip"`
	Standard     string
	Ports        []Port
	Architecture string
	Messenger    string
	Rating       int
	Status       string
}

type Port struct {
	ID       uint   `gorm:"primary_key"`
	Number   int    `gorm:"unique_index:idx_port"`
	Protocol string `gorm:"unique_index:idx_port"`
	Status   string `gorm:"unique_index:idx_port"`
	PeerID   uint   `gorm:"unique_index:idx_port"`
	Peer     *Peer
}

type Service struct {
	ID      uint   `gorm:"primary_key"`
	Name    string `gorm:"unique_index:idx_service"`
	Version string
	Product string
	OsType  string
	PortID  uint `gorm:"unique_index:idx_service"`
	Port    *Port
}

type Scan struct {
	Name      string
	Target    string
	Status    int
	Outfolder string
	Outfile   string
	Cmd       string
}

func (s Step) String() string {
	return [...]string{"NOT_DEFINED", "IMPORTED", "SWEEPED", "NEW", "SCANNED"}[s]
}

var log = logger.GetOrCreate("vat")

/*=====
General
======*/
func InitDB(dbpath string) *gorm.DB {
	// Create connection to DB
	db, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		fmt.Println(fmt.Sprintf("[DB ERROR] %s", err))
		os.Exit(1)
	}
	// Disable logging
	if os.Getenv("DEBUG") == "1" {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	// Migrate schema
	migrateDB(db)

	return db
}

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&Service{})
	db.AutoMigrate(&Port{})
	db.AutoMigrate(&Peer{})
}

/*==========
Gatherers
============*/

func (s *Service) GetPort(db *gorm.DB) *Port {
	port := &Port{}
	db.Where("id = ?", s.PortID).Find(&port)
	return port
}

func (p *Port) GetPeer(db *gorm.DB) *Peer {
	peer := &Peer{}
	db.Where("id = ?", p.PeerID).Find(&peer)
	return peer
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

/*==========
// Constructors
============*/
func AddPeer(db *gorm.DB, id uint, protocol string,
	standard string,
	architecture string,
	messenger string,
	address string,
	rating int,
	status string) *Peer {
	lock.Lock()
	defer lock.Unlock()

	t := &Peer{
		ID:           id,
		Protocol:     protocol,
		Standard:     standard,
		Architecture: architecture,
		Messenger:    messenger,
		Address:      address,
		Rating:       rating,
		Status:       status,
	}
	db.Create(t)
	return t
}

func AddPort(db *gorm.DB, number int, protocol, status string, h *Peer) (*Port, bool) {
	lock.Lock()
	defer lock.Unlock()

	duplicate := false
	t := &Port{
		Number:   number,
		Protocol: protocol,
		Status:   status,
		Peer:     h,
	}
	if err := db.Create(t).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			duplicate = true
		}
	}

	return t, duplicate
}

func AddService(db *gorm.DB, name, version, product, osType string, p *Port, pID uint) *Service {
	lock.Lock()
	defer lock.Unlock()

	t := &Service{
		Name:    name,
		Version: version,
		Product: product,
		OsType:  osType,
		Port:    p,
		PortID:  pID,
	}
	db.Create(t)
	return t
}
