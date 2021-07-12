package mock

import (
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
)

// PeerAuthenticationProcessorStub -
type PeerAuthenticationProcessorStub struct {
	ProcessReceivedCalled func(message p2p.MessageP2P, peerHeartbeat process.InterceptedPeerAuthentication) error
}

// ProcessReceived -
func (proc *PeerAuthenticationProcessorStub) ProcessReceived(message p2p.MessageP2P, peerHeartbeat process.InterceptedPeerAuthentication) error {
	if proc.ProcessReceivedCalled != nil {
		return proc.ProcessReceivedCalled(message, peerHeartbeat)
	}

	return nil
}

// IsInterfaceNil -
func (proc *PeerAuthenticationProcessorStub) IsInterfaceNil() bool {
	return proc == nil
}
