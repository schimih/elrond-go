package evaluation

import (
	"github.com/elrond-go/cmd/vat/scan"
	"github.com/elrond-go/cmd/vat/utils"
)

type identity struct {
	address string
	ports   []scan.Port
	status  string
}

func newIdentity(address string, ports []scan.Port) (id identity) {
	return identity{
		address: address,
		ports:   ports,
		status:  string(utils.NEW),
	}
}
