package disabled

import (
	"errors"

	"github.com/ElrondNetwork/elrond-go/data/api"
)

var errCannotReturnAccountListFromMetachainNode = errors.New("account list cannot be returned by a metachain")

type accountListProcessor struct{}

// NewDisabledAccountListProcessor returns a disabled implementation to be used on shard nodes
func NewDisabledAccountListProcessor() *accountListProcessor {
	return &accountListProcessor{}
}

// GetAccountsList returns the errCannotReturnAccountListFromMetachainNode error
func (dlp *accountListProcessor) GetAccountsList() ([]*api.Account, error) {
	return nil, errCannotReturnAccountListFromMetachainNode
}

// IsInterfaceNil returns true if there is no value under the interface
func (dlp *accountListProcessor) IsInterfaceNil() bool {
	return dlp == nil
}
