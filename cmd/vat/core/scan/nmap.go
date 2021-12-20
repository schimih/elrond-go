package scan

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/VAT/cmd/vul/core/model"
	"github.com/VAT/cmd/vul/core/utils"
	go_nmap "github.com/lair-framework/go-nmap"
)

type NmapScan model.Scan

// Constructor for NmapScan
func NewScan(name, target, folder, file, nmapArgs string) *NmapScan {
	// Create a Scan
	s := &NmapScan{
		Name:   name,
		Target: target,
		Status: model.NOT_STARTED,
	}
	// Construct output path and create if it doesn't exist
	s.Outfolder = filepath.Join(utils.Config.Outfolder, utils.CleanPath(target), folder)
	s.Outfile = filepath.Join(s.Outfolder, utils.CleanPath(file))
	utils.CheckDir(s.Outfolder)
	// Construct command
	s.Cmd = s.constructCmd(nmapArgs)
	return s
}

func (s *NmapScan) preScan() {
	s.Status = model.IN_PROGRESS
}
func (s *NmapScan) postScan() {
	s.Status = model.FINISHED
}

func (s *NmapScan) constructCmd(args string) string {
	// TODO: relative path
	// strutocamila
	usr, _ := user.Current()
	usrName := strings.SplitN(usr.HomeDir, "\\", 3)
	return fmt.Sprintf("nmap %s %s -oA C:/Users/%s/vat/%s/%s_%s", args, s.Target, usrName[2], s.Target, s.Name, s.Target)
}

// Run nmap scan
func (s *NmapScan) RunNmap() {
	// Pre-scan checks
	s.preScan()
	// Run nmap
	_, err := utils.ShellCmd(s.Cmd)
	if err != nil {
		s.Status = model.FAILED
	}
	// Post-scan checks
	s.postScan()
}

// Parse nmap XML output file
func (s *NmapScan) ParseOutput() *go_nmap.NmapRun {
	// ???? check for cross platform compatibility
	sweepXML := fmt.Sprintf("%s.xml", s.Outfile) //s.Outfile
	dat, err := ioutil.ReadFile(sweepXML)
	if err != nil {
		log.Info(fmt.Sprintf("Error while opening output file: %s", sweepXML))
		return nil
	}

	res, err := go_nmap.Parse(dat)
	if err != nil {
		log.Info("Error while parsing nmap output")
		return nil
	}
	return res
}
