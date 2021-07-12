package mock

import (
	"sync"

	"github.com/ElrondNetwork/elrond-go/core"
)

// InterceptedPeerAuthenticationStub -
type InterceptedPeerAuthenticationStub struct {
	InterceptedDataStub
	PublicKeyCalled    func() []byte
	PeerIDCalled       func() core.PeerID
	SignatureCalled    func() []byte
	mutComputedShardID sync.Mutex
	computedShardID    uint32
}

// PeerID -
func (i *InterceptedPeerAuthenticationStub) PeerID() core.PeerID {
	if i.PeerIDCalled != nil {
		return i.PeerIDCalled()
	}

	return ""
}

// Signature -
func (i *InterceptedPeerAuthenticationStub) Signature() []byte {
	if i.SignatureCalled != nil {
		return i.SignatureCalled()
	}

	return make([]byte, 0)
}

// ComputedShardID -
func (i *InterceptedPeerAuthenticationStub) ComputedShardID() uint32 {
	i.mutComputedShardID.Lock()
	defer i.mutComputedShardID.Unlock()

	return i.computedShardID
}

// PublicKey -
func (i *InterceptedPeerAuthenticationStub) PublicKey() []byte {
	if i.PublicKeyCalled != nil {
		return i.PublicKeyCalled()
	}

	return make([]byte, 0)
}

// SetComputedShardID -
func (i *InterceptedPeerAuthenticationStub) SetComputedShardID(shardId uint32) {
	i.mutComputedShardID.Lock()
	i.computedShardID = shardId
	i.mutComputedShardID.Unlock()
}
