package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/VAT/cmd/vul/core/model"
	"github.com/jinzhu/gorm"
)

// UTILS
// General features, should continuously be updated.
var Config config
var Const_notification_delay_unit = 10

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

//config for db
type config struct {
	Outfolder string
	DB        *gorm.DB
	DBPath    string
}

var log = logger.GetOrCreate("vat")

// Initialize global config (db, etc.)
// From now on it will be accessible as utils.Config
func InitConfig() {
	Config = config{}
	t := time.Now()
	// Create output folder
	if os.Getenv("OUT_FOLDER") != "" {
		Config.Outfolder = filepath.Join(os.Getenv("OUT_FOLDER"), "vat")
	} else {
		usr, _ := user.Current()
		Config.Outfolder = filepath.Join(usr.HomeDir, "vat")
	}
	CheckDir(Config.Outfolder)

	// Init DB
	if os.Getenv("VAT_DB_PATH") != "" {
		Config.DBPath = os.Getenv("VAT_DB_PATH")
	} else {
		Config.DBPath = filepath.Join(Config.Outfolder, "vat"+t.Format("2006-01-02")+".db")
	}
	Config.DB = model.InitDB(Config.DBPath)
	log.Info("DB ready to be used:", "DB", Config.DBPath)
}

// Check directory existance, or create it otherwise
func CheckDir(dir string) {
	// Create a directory if doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		log.Info("Created directory:", "path", dir)
	}
}

func ShellCmd(cmd string) (string, error) {
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			log.Error(err.Error())
		}
		return string(output), err
	}
	return string(output), err
}

// Replace slashes with underscores, when the string is used in a path
func CleanPath(s string) string {
	return strings.Replace(s, "/", "_", -1)
}

//loadFile from json file
func LoadFile(src string) bool {
	// If it's a folder, iterate through all the files contained in there
	//we should initialize the peers array
	var peers []model.Peer
	fpath, _ := os.Stat(src)
	if fpath.IsDir() {
		err := filepath.Walk(src,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Error(fmt.Sprintf("Error while listing content of directory: %s", src))
					return err
				}
				t, _ := os.Stat(path)
				if filepath.Ext(t.Name()) == ".json" {
					//when there will be time, movee this into a separate function
					jsonFile, err := os.Open(path)
					if err != nil {
						return err
					}
					// read our opened jsonFile as a byte array.
					byteValue, _ := ioutil.ReadAll(jsonFile)
					//in peers we will have all results from json file
					json.Unmarshal(byteValue, &peers)
					//move to db and use it for operational scope
					for _, peer := range peers {
						model.AddPeer(Config.DB, peer.ID, peer.Protocol,
							peer.Standard, peer.Architecture, peer.Messenger,
							peer.Address, peer.Rating, model.IMPORTED.String())
					}
				}
				return nil
			})
		if err != nil {
			return false
		}
	} else {
		// If it's a file, import it straight away
		if filepath.Ext(fpath.Name()) != ".json" {
			log.Error(fmt.Sprintf("Please provide a .json file"))
			return false
		}
		//when there will be time, move this into a separate function
		jsonFile, err := os.Open(src)
		if err != nil {
			return false
		}
		// read our opened jsonFile as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)
		//in peers we will have all results from json file
		json.Unmarshal(byteValue, &peers)
		for _, peer := range peers {
			model.AddPeer(Config.DB, peer.ID, peer.Protocol,
				peer.Standard, peer.Architecture, peer.Messenger,
				peer.Address, peer.Rating, model.IMPORTED.String())
		}
	}
	return true
}

func ScanShow(item string) {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := model.GetAllPeers(Config.DB)
	if len(peersDB) == 0 {
		log.Info("No peers in DB. First load a json or run discovery!")
		return
	}
	dataLines := make([]*display.LineData, 0, len(peersDB))

	for idx, p := range peersDB {
		rAddress := p.Address
		for jdx, tPort := range p.GetPorts(Config.DB) {
			rPort := fmt.Sprintf("%d", tPort.Number)
			rStatus := tPort.Status
			rProtocol := tPort.Protocol
			rIndex := fmt.Sprintf("%d", idx)
			horizontalLineAfter := jdx == len(p.GetPorts(Config.DB))-1
			lines := display.NewLineData(horizontalLineAfter, []string{rIndex, rAddress, rPort, rStatus, rProtocol})
			dataLines = append(dataLines, lines)
		}
	}

	table, _ := display.CreateTableString(header, dataLines)
	fmt.Println(table)
}
