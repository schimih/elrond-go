package mock

import (
	"github.com/ElrondNetwork/elrond-go/process"
)

// PeerAuthenticationProcessorStub -
type PeerAuthenticationProcessorStub struct {
	ProcessReceivedCalled  func(originalPayload []byte, peerHeartbeat process.InterceptedPeerAuthentication) error
	IsHardforkMessageCaled func(peerHeartbeat process.InterceptedPeerAuthentication) bool
}

// ProcessReceived -
func (proc *PeerAuthenticationProcessorStub) ProcessReceived(originalPayload []byte, peerHeartbeat process.InterceptedPeerAuthentication) error {
	if proc.ProcessReceivedCalled != nil {
		return proc.ProcessReceivedCalled(originalPayload, peerHeartbeat)
	}

	return nil
}

// IsHardforkMessage -
func (proc *PeerAuthenticationProcessorStub) IsHardforkMessage(peerHeartbeat process.InterceptedPeerAuthentication) bool {
	if proc.IsHardforkMessageCaled != nil {
		return proc.IsHardforkMessageCaled(peerHeartbeat)
	}

	return false
}

// IsInterfaceNil -
func (proc *PeerAuthenticationProcessorStub) IsInterfaceNil() bool {
	return proc == nil
}
