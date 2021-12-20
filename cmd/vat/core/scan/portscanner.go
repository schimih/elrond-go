package scan

import (
	"fmt"
	"sync"
	"time"

	"github.com/VAT/cmd/vul/core/model"
	"github.com/VAT/cmd/vul/core/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

/*=====
Consts
=====*/
var ScansList = []*NmapScan{}
var notificationDelay time.Duration = time.Duration(utils.Const_notification_delay_unit) * time.Second

/*=====
Scanner
=====*/
func ScanPort(scanType string) {
	/*
			for later
		   1. Check for the existence of the .json
		   2. if it does not exist -> inform the user to run discovery first
		   3. load data from the .json
			first
		   4. check the type of the check (elrond port/80:8080/22)
		   5. run the check
		   6. create a new file -> store results of the port check + new rating and reason for the rating
	*/

	folder := "/"
	// Dispatch scan
	switch scanType {
	case "TCP-ELROND":
		log.Info("Verifying if just the elrond port list (37373-38383) is open")
		file, nmapArgs := "tcp_elrond", utils.NMAP_TCP_ELROND
		startPortScan(file, folder, file, nmapArgs)
	case "TCP-WEB":
		log.Info("Verifying if port 80 or 8080 is accessible")
		file, nmapArgs := "tcp_web", utils.NMAP_TCP_WEB
		startPortScan(file, folder, file, nmapArgs)
	case "TCP-SSH":
		log.Info("Verifying if the ssh port 22 is accessible")
		file, nmapArgs := "tcp_ssh", utils.NMAP_TCP_SSH
		startPortScan(file, folder, file, nmapArgs)
	case "TCP-ALL":
		log.Info("Scans everything")
		file, nmapArgs := "tcp_all", utils.NMAP_TCP_STANDARD
		startPortScan(file, folder, file, nmapArgs)
	case "TCP-LOAD":
		log.Info("Loading an analysis")
		file, nmapArgs := "tcp_load", utils.NMAP_TCP_ELROND
		startPortScan(file, folder, file, nmapArgs)
	default:
		log.Error("%s is an invalid command for portscan", scanType)
		return
	}

}

func startPortScan(name, folder, file, nmapArgs string) {
	var wg sync.WaitGroup //=
	peers := model.GetAllPeers(utils.Config.DB)
	for _, h := range peers {
		wg.Add(1) //=

		go func() { //=
			defer wg.Done() //=
			temp := h
			fname := fmt.Sprintf("%s_%s", file, h.Address)
			worker(name, &temp, folder, fname, nmapArgs)
		}()
		wg.Wait() //=
		log.Info("Scanning done for peer:", "address", h.Address)

	}

}

func worker(name string, h *model.Peer, folder string, file string, nmapArgs string) {
	// Instantiate new NmapScan
	s := NewScan(name, h.Address, folder, file, nmapArgs)
	ScansList = append(ScansList, s)
	// Run the scan
	s.RunNmap()
	// Parse nmap's output
	res := s.ParseOutput()
	if res != nil {
		// hosts - because nmap
		for _, record := range res.Hosts {
			ProcessResults(h, record)
		}
	}
}

func ProcessResults(peer *model.Peer, record go_nmap.Host) {
	// Parse ports
	for _, port := range record.Ports {
		// Create new port, will add to db if new
		np, _ := model.AddPort(utils.Config.DB, port.PortId, port.Protocol, port.State.State, peer)

		// Add Service
		if port.Service.Name != "" {
			model.AddService(utils.Config.DB, port.Service.Name, port.Service.Version, port.Service.Product, port.Service.OsType, np, np.ID)
		}
	}

	model.Mutex.Lock()
	peer.Status = model.SCANNED.String()
	utils.Config.DB.Save(&peer)
	model.Mutex.Unlock()
}
