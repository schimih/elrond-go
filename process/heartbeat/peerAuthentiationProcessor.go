package heartbeat

import (
	"bytes"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/process"
)

// PeerAuthenticationProcessorArgs is the argument for the NewPeerAuthenticationProcessor function
type PeerAuthenticationProcessorArgs struct {
	PeerSignatureHandler     crypto.PeerSignatureHandler
	NetworkShardingCollector process.NetworkShardingCollector
	HardforkTrigger          process.HardforkTrigger
}

type peerAuthenticationProcessor struct {
	peerSignatureHandler     crypto.PeerSignatureHandler
	networkShardingCollector process.NetworkShardingCollector
	hardforkTrigger          process.HardforkTrigger
}

// NewPeerAuthenticationProcessor creates a new peer authentication processor
func NewPeerAuthenticationProcessor(args PeerAuthenticationProcessorArgs) (*peerAuthenticationProcessor, error) {
	if check.IfNil(args.PeerSignatureHandler) {
		return nil, process.ErrNilPeerSignatureHandler
	}
	if check.IfNil(args.HardforkTrigger) {
		return nil, process.ErrNilHardforkTrigger
	}
	if check.IfNil(args.NetworkShardingCollector) {
		return nil, process.ErrNilNetworkShardingCollector
	}
}

// ProcessReceived processes the received message
func (pap *peerAuthenticationProcessor) ProcessReceived(originalPayload []byte, peerHeartbeat process.InterceptedPeerAuthentication) error {
	if check.IfNil(peerHeartbeat) {
		return process.ErrNilPeerAuthenticationInterceptedData
	}

	err := pap.peerSignatureHandler.VerifyPeerSignature(peerHeartbeat.PublicKey(), peerHeartbeat.PeerID(), peerHeartbeat.Signature())
	if err != nil {
		return err
	}

	pap.networkShardingCollector.UpdatePeerIDInfo(peerHeartbeat.PeerID(), peerHeartbeat.PublicKey(), peerHeartbeat.ComputedShardID())

	isHardforkTrigger, err := pap.hardforkTrigger.TriggerReceived(originalPayload, peerHeartbeat.HardforkPayload(), peerHeartbeat.PublicKey())
	if isHardforkTrigger {
		return err
	}

	return nil
}

// IsHardforkMessage will return true if the peerHeartbeat instance is a correctly formatted hardfork message from the expected peer
func (pap *peerAuthenticationProcessor) IsHardforkMessage(peerHeartbeat process.InterceptedPeerAuthentication) bool {
	if check.IfNil(peerHeartbeat) {
		return false
	}

	return bytes.Equal(peerHeartbeat.PublicKey(), pap.hardforkTrigger.HardforkTriggerPublicKeyBytes())
}

// IsInterfaceNil returns true if there is no value under the interface
func (pap *peerAuthenticationProcessor) IsInterfaceNil() bool {
	return pap == nil
}
