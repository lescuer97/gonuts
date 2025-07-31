package nostr

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	n "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/keyer"
	"github.com/nbd-wtf/go-nostr/nip17"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type NostrRelayConnections struct {
	url  string
	conn *n.Connection
}

type NostrClient struct {
	relaysUrl  []string
	simplePool *n.SimplePool
	keyer      n.Keyer
	chanEvent  chan n.Event
}

var (
	DefaultRelays = []string{"wss://relay.damus.io/", "wss://relay.primal.net/", "wss://relay.snort.social/"}
)

func SetupNostrClient(privateKey *secp256k1.PrivateKey, timestamp n.Timestamp, relays []string) (*NostrClient, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("passed private key is nil")
	}
	ctx := context.Background()

	selectedRelays := relays

	sec, err := nip19.EncodePrivateKey(hex.EncodeToString(privateKey.Serialize()))
	defer func() { sec = "" }()
	if err != nil {
		return nil, err
	}

	keyer, err := keyer.NewPlainKeySigner(sec)
	if err != nil {
		return nil, err
	}

	simplePool := n.NewSimplePool(ctx)
	if simplePool == nil {
		return nil, fmt.Errorf("Simple pool is nil")
	}

	client := NostrClient{
		relaysUrl:  selectedRelays,
		keyer:      keyer,
		simplePool: simplePool,
	}

	eventWatch := nip17.ListenForMessages(ctx, client.simplePool, client.keyer, client.relaysUrl, timestamp)

	client.chanEvent = eventWatch

	return &client, nil
}

// Close properly cleans up the NostrClient resources
// This should be called when the client is no longer needed to prevent memory leaks
func (nc *NostrClient) Close() error {
	if nc == nil {
		return nil
	}

	// Close the event channel if it exists
	if nc.chanEvent != nil {
		close(nc.chanEvent)
		nc.chanEvent = nil
	}

	// Close the simple pool to terminate all relay connections
	if nc.simplePool != nil {
		for _, url := range nc.relaysUrl {
			nc.simplePool.Close(url)
		}
		nc.simplePool = nil
	}

	// Clear sensitive keyer data for security
	nc.keyer = nil

	// Clear other fields
	nc.relaysUrl = nil

	return nil
}

// GetKeyer returns the keyer for the NostrClient
func (nc *NostrClient) GetKeyer() n.Keyer {
	return nc.keyer
}

// GetEventsChan returns the event channel for the NostrClient
func (nc *NostrClient) GetEventsChan() chan n.Event {
	return nc.chanEvent
}
