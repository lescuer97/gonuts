package wallet

import (
	"encoding/json"
	"log"

	"github.com/elnosh/gonuts/cashu/nuts/nut18"
	n "github.com/nbd-wtf/go-nostr"
)

// handleNostrEvent processes a single nostr event.
func (w *Wallet) handleNostrEvent(event n.Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in nostr event processing: %v", r)
		}
	}()

	var payload nut18.PaymentRequestPayload
	err := json.Unmarshal([]byte(event.Content), &payload)
	if err != nil {
		// not a payment request, ignore
		return
	}

	_, err = w.ReceivePaymentRequestPayload(payload)
	if err != nil {
		log.Printf("error processing payment request payload: %v", err)
	}
}

func processNostrRequests(wallet *Wallet) {
	if wallet.nostrClient == nil || wallet.nostrClient.GetEventsChan() == nil {
		log.Println("Nostr client or event channel is not available, stopping request processing.")
		return
	}
	eventsChan := wallet.nostrClient.GetEventsChan()

	for event := range eventsChan {
		go wallet.handleNostrEvent(event)
	}
}
