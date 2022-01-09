package scan

import (
	go_nmap "github.com/lair-framework/go-nmap"
)

type Scanner interface {
	Scan() *go_nmap.NmapRun
	IsInterfaceNil() bool
}
