// Based on perun-examples/payment-channel/client/handle.go (Apache-2.0, PolyCrypt GmbH).
// Changes for the week 5 labs:
//   - Added `return` after Reject. Set c.Buggy = true to get the original buggy behavior back.
//   - Funding is checked against c.Policy.MaxPeerFunding instead of requiring 0.
//   - Payments smaller than c.Policy.MinPayment are rejected.
//   - Each adjudicator event type is logged in detail.

package client

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// HandleProposal is the callback for incoming channel proposals.
func (c *PaymentClient) HandleProposal(p client.ChannelProposal, r *client.ProposalResponder) {
	lcp, err := func() (*client.LedgerChannelProposalMsg, error) {
		// Ensure that we got a ledger channel proposal.
		lcp, ok := p.(*client.LedgerChannelProposalMsg)
		if !ok {
			return nil, fmt.Errorf("invalid proposal type: %T", p)
		}

		// Check that we have the correct number of participants.
		if lcp.NumPeers() != 2 {
			return nil, fmt.Errorf("invalid number of participants: %d", lcp.NumPeers())
		}

		// Check that the channel has the expected assets and funding balances.
		const assetIdx, peerIdx = 0, 1
		if err := channel.AssertAssetsEqual(lcp.InitBals.Assets, []channel.Asset{c.currency}); err != nil {
			return nil, fmt.Errorf("invalid assets: %v", err)
		}
		if lcp.FundingAgreement[assetIdx][peerIdx].Cmp(c.Policy.MaxPeerFunding) > 0 {
			return nil, fmt.Errorf("peer funding too high: %v ETH > %v ETH",
				WeiToEth(lcp.FundingAgreement[assetIdx][peerIdx]), WeiToEth(c.Policy.MaxPeerFunding))
		}
		return lcp, nil
	}()
	if err != nil {
		log.Printf("[%s] REJECT proposal: %v", c.Name, err)
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
		if !c.Buggy {
			return // The fix. The original code lacks this line and falls through to lcp.Accept with lcp == nil.
		}
	}

	// Create a channel accept message and send it.
	accept := lcp.Accept(
		c.account,                // The account we use in the channel.
		client.WithRandomNonce(), // Our share of the channel nonce.
	)
	ch, err := r.Accept(context.TODO(), accept)
	if err != nil {
		log.Printf("[%s] error accepting channel proposal: %v", c.Name, err)
		return
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	// Store channel.
	c.channels <- newPaymentChannel(ch, c.currency)
}

// HandleUpdate is the callback for incoming channel updates.
func (c *PaymentClient) HandleUpdate(cur *channel.State, next client.ChannelUpdate, r *client.UpdateResponder) {
	// We accept every update that increases our balance.
	err := func() error {
		err := channel.AssertAssetsEqual(cur.Assets, next.State.Assets)
		if err != nil {
			return fmt.Errorf("invalid assets: %v", err)
		}

		receiverIdx := 1 - next.ActorIdx // This works because we are in a two-party channel.
		curBal := cur.Allocation.Balance(receiverIdx, c.currency)
		nextBal := next.State.Allocation.Balance(receiverIdx, c.currency)
		if nextBal.Cmp(curBal) < 0 {
			return fmt.Errorf("invalid balance: %v", nextBal)
		}

		// Anti-spam rule: reject if we receive more than 0 but less than MinPayment.
		// The finalize update (IsFinal) transfers 0, so it still passes.
		received := new(big.Int).Sub(nextBal, curBal)
		if c.Policy.MinPayment != nil && received.Sign() > 0 && received.Cmp(c.Policy.MinPayment) < 0 {
			return fmt.Errorf("payment too small: %v ETH < %v ETH", WeiToEth(received), WeiToEth(c.Policy.MinPayment))
		}
		return nil
	}()
	if err != nil {
		log.Printf("[%s] REJECT update v%d: %v", c.Name, next.State.Version, err)
		r.Reject(context.TODO(), err.Error()) //nolint:errcheck // It's OK if rejection fails.
		if !c.Buggy {
			return // The fix. The original code lacks this line, so it calls r.Accept and panics.
		}
	}

	// Send the acceptance message.
	err = r.Accept(context.TODO())
	if err != nil {
		panic(err)
	}
}

// HandleAdjudicatorEvent is the callback for smart contract events.
func (c *PaymentClient) HandleAdjudicatorEvent(e channel.AdjudicatorEvent) {
	switch ev := e.(type) {
	case *channel.RegisteredEvent:
		log.Printf("[%s] EVENT Registered  version=%d timeout=%v", c.Name, ev.State.Version, ev.Timeout())
	case *channel.ProgressedEvent:
		log.Printf("[%s] EVENT Progressed  version=%d by idx=%d", c.Name, ev.State.Version, ev.Idx)
	case *channel.ConcludedEvent:
		log.Printf("[%s] EVENT Concluded   version=%d", c.Name, ev.Version())
	default:
		log.Printf("[%s] EVENT %T", c.Name, e)
	}
}
