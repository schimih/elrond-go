package factory

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/node/external"
	"github.com/ElrondNetwork/elrond-go/node/trieIterators"
	"github.com/ElrondNetwork/elrond-go/node/trieIterators/disabled"
)

// CreateAccountsListHandler will create a new instance of AccountsListHandler
func CreateAccountsListHandler(args trieIterators.ArgTrieIteratorProcessor) (external.AccountsListHandler, error) {
	//TODO add unit tests
	if args.ShardID == core.MetachainShardId {
		return disabled.NewDisabledAccountListProcessor(), nil
	}

	return trieIterators.NewAccountListProcessor(args)
}
