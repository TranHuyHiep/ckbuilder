// Based on perun-examples/payment-channel/client/channel.go (Apache-2.0, PolyCrypt GmbH).
// Changes for the week 5 labs: every function returns an error instead of panicking,
// plus inspection helpers (Balances, Version) and attack / dispute helpers.

package client

import (
	"context"
	"math/big"
	"time"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// PaymentChannel is a wrapper for a Perun channel for the payment use case.
type PaymentChannel struct {
	ch       *client.Channel
	currency channel.Asset
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel, currency channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:       ch,
		currency: currency,
	}
}

// SendPayment sends a payment to the channel peer.
func (c PaymentChannel) SendPayment(amount float64) error {
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		ethAmount := EthToWei(big.NewFloat(amount))
		actor := c.ch.Idx()
		peer := 1 - actor
		state.Allocation.TransferBalance(actor, peer, c.currency, ethAmount)
	})
}

// Balances returns the current in-channel balances of [index 0, index 1] in ETH.
func (c PaymentChannel) Balances() []*big.Float {
	s := c.ch.State()
	return []*big.Float{
		WeiToEth(s.Allocation.Balance(0, c.currency)),
		WeiToEth(s.Allocation.Balance(1, c.currency)),
	}
}

// Version returns the version of the current state.
func (c PaymentChannel) Version() uint64 {
	return c.ch.State().Version
}

// IsFinal reports whether the current state is final.
func (c PaymentChannel) IsFinal() bool {
	return c.ch.State().IsFinal
}

// Finalize sets IsFinal = true with an off-chain update (needs the peer's consent).
func (c PaymentChannel) Finalize() error {
	if c.ch.State().IsFinal {
		return nil
	}
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		state.IsFinal = true
	})
}

// Settle finalizes (if needed), then concludes and withdraws: the fast path.
func (c PaymentChannel) Settle() error {
	if err := c.Finalize(); err != nil {
		return err
	}
	if err := c.ch.Settle(context.TODO(), false); err != nil {
		return err
	}
	return c.ch.Close()
}

// ForceSettle settles without finalizing, i.e. without the peer's consent:
// register the state on-chain, wait out the challenge duration, conclude, withdraw.
func (c PaymentChannel) ForceSettle() (time.Duration, error) {
	start := time.Now()
	if err := c.ch.Settle(context.TODO(), false); err != nil {
		return time.Since(start), err
	}
	d := time.Since(start)
	return d, c.ch.Close()
}

// WithdrawOnly only withdraws, assuming the other side already concluded the channel.
func (c PaymentChannel) WithdrawOnly() error {
	if err := c.ch.Settle(context.TODO(), true); err != nil {
		return err
	}
	return c.ch.Close()
}

// MintForMe tries to raise our own balance without debiting anyone (attack, Lab 03).
func (c PaymentChannel) MintForMe(amount float64) error {
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		me := c.ch.Idx()
		bal := state.Allocation.Balance(me, c.currency)
		newBal := new(big.Int).Add(bal, EthToWei(big.NewFloat(amount)))
		state.Allocation.SetBalance(me, c.currency, newBal)
	})
}

// Steal moves funds from the peer to us (attack, Lab 03).
func (c PaymentChannel) Steal(amount float64) error {
	return c.ch.Update(context.TODO(), func(state *channel.State) {
		me := c.ch.Idx()
		state.Allocation.TransferBalance(1-me, me, c.currency, EthToWei(big.NewFloat(amount)))
	})
}
