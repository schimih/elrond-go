package scan

import (
	go_nmap "github.com/lair-framework/go-nmap"
)

type NmapWrapper interface {
	RunNmap()
	ParseOutput() *go_nmap.NmapRun
	IsListNil() bool
}
