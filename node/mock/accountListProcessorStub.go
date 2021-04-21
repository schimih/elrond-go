package mock

import (
	"github.com/ElrondNetwork/elrond-go/data/api"
)

// AccountListProcessorStub -
type AccountListProcessorStub struct {
	GetAccountsListCalled      func() ([]*api.Account, error)
}

// GetAccountsList -
func (als *AccountListProcessorStub) GetAccountsList() ([]*api.Account, error) {
	if als.GetAccountsListCalled != nil {
		return als.GetAccountsListCalled()
	}

	return nil, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (als *AccountListProcessorStub) IsInterfaceNil() bool {
	return als == nil
}