package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/ElrondNetwork/elrond-go-core/display"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/VAT/cmd/vul/core/result"
)

var log = logger.GetOrCreate("vat")

// Check directory existance, or create it otherwise
func CheckDir(dir string) {
	// Create a directory if doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		log.Info("Created directory:", "path", dir)
	}
}

// Replace slashes with underscores, when the string is used in a path
func CleanPath(s string) string {
	return strings.Replace(s, "/", "_", -1)
}

//loadFile from json file
func LoadFile(src string) bool {
	// If it's a folder, iterate through all the files contained in there
	//we should initialize the peers array
	// var peers []result.Peer
	// fpath, _ := os.Stat(src)
	// if fpath.IsDir() {
	// 	err := filepath.Walk(src,
	// 		func(path string, info os.FileInfo, err error) error {
	// 			if err != nil {
	// 				log.Error(fmt.Sprintf("Error while listing content of directory: %s", src))
	// 				return err
	// 			}
	// 			t, _ := os.Stat(path)
	// 			if filepath.Ext(t.Name()) == ".json" {
	// 				//when there will be time, movee this into a separate function
	// 				jsonFile, err := os.Open(path)
	// 				if err != nil {
	// 					return err
	// 				}
	// 				// read our opened jsonFile as a byte array.
	// 				byteValue, _ := ioutil.ReadAll(jsonFile)
	// 				//in peers we will have all results from json file
	// 				json.Unmarshal(byteValue, &peers)
	// 				//move to db and use it for operational scope
	// 				for _, peer := range peers {
	// 					result.AddPeer(Config.DB, peer.ID, peer.Protocol,
	// 						peer.Standard, peer.Architecture, peer.Messenger,
	// 						peer.Address, peer.Rating, result.IMPORTED.String())
	// 				}
	// 			}
	// 			return nil
	// 		})
	// 	if err != nil {
	// 		return false
	// 	}
	// } else {
	// 	// If it's a file, import it straight away
	// 	if filepath.Ext(fpath.Name()) != ".json" {
	// 		log.Error(fmt.Sprintf("Please provide a .json file"))
	// 		return false
	// 	}
	// 	//when there will be time, move this into a separate function
	// 	jsonFile, err := os.Open(src)
	// 	if err != nil {
	// 		return false
	// 	}
	// 	// read our opened jsonFile as a byte array.
	// 	byteValue, _ := ioutil.ReadAll(jsonFile)
	// 	//in peers we will have all results from json file
	// 	json.Unmarshal(byteValue, &peers)
	// 	for _, peer := range peers {
	// 		result.AddPeer(Config.DB, peer.ID, peer.Protocol,
	// 			peer.Standard, peer.Architecture, peer.Messenger,
	// 			peer.Address, peer.Rating, result.IMPORTED.String())
	// 	}
	// }
	return true
}

func ScanShow(item string) {
	header := []string{"Index", "Address", "Port", "Status", "Service"}
	peersDB := result.GetAllPeers(Config.DB)
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
