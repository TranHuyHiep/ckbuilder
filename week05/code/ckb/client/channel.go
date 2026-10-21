// Copyright 2024 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"log"
	"math/big"
	"time"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/perun-ckb-backend/channel/asset"
)

type PaymentChannel struct {
	ch     *client.Channel
	assets []channel.Asset
}

// newPaymentChannel creates a new payment channel.
func newPaymentChannel(ch *client.Channel, assets []channel.Asset) *PaymentChannel {
	return &PaymentChannel{
		ch:     ch,
		assets: assets,
	}
}

func (c PaymentChannel) State() *channel.State {
	return c.ch.State().Clone()
}

func (c PaymentChannel) SendPayment(amounts map[channel.Asset]float64) {
	// Transfer the given amount from us to peer.
	// Use UpdateBy to update the channel state.
	err := c.ch.Update(context.TODO(), func(state *channel.State) {
		actor := c.ch.Idx()
		peer := 1 - actor
		for a, amount := range amounts {
			if amount < 0 {
				continue
			}
			switch a := a.(type) {
			case *asset.Asset:
				if a.IsCKBytes {
					shannonAmount := CKByteToShannon(big.NewFloat(amount))
					state.Allocation.TransferBalance(actor, peer, a, shannonAmount)
				} else {
					sudtAmmount := new(big.Int).SetUint64(uint64(amount))
					state.Allocation.TransferBalance(actor, peer, a, sudtAmmount)
				}
			default:
				panic("Asset is not of type *asset.Asset")
			}
		}

	})
	if err != nil {
		panic(err)
	}

}

// Settle settles the payment channel and withdraws the funds.
func (c PaymentChannel) Settle() {
	// Finalize the channel to enable fast settlement.
	if !c.ch.State().IsFinal {
		err := c.ch.Update(context.TODO(), func(state *channel.State) {
			state.IsFinal = true
		})
		if err != nil {
			panic(err)
		}
	}

	// Settle concludes the channel and withdraws the funds.
	err := c.ch.Settle(context.TODO(), false)
	if err != nil {
		panic(err)
	}

	// Close frees up channel resources.
	c.ch.Close()
}

// ForceSettle settles without finalizing (Lab 05, like Lab 04): puts the latest
// state on-chain with a Dispute transaction, waits out the challenge duration, then ForceClose.
//
// Workaround for a race in perun-ckb-backend: go-perun waits for a Registered event, but the
// backend's polling subscription first polls 4 s after it is created. On a fast devnet the
// Dispute transaction is already committed by then, the subscription takes the disputed status
// as its starting point and never emits Registered, so Settle would wait forever. The channel
// phase is already Registered at that point, so we bound the first Settle with a timeout and
// call Settle again: ensureRegistered then returns immediately and ForceClose runs.
func (c PaymentChannel) ForceSettle() {
	start := time.Now()
	wait := time.Duration(c.ch.Params().ChallengeDuration)*time.Second + 20*time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	err := c.ch.Settle(ctx, false)
	cancel()
	if err != nil {
		if c.ch.Phase() != channel.Registered {
			panic(err)
		}
		log.Printf("Dispute is on-chain (phase %v) but no Registered event arrived within %v "+
			"(perun-ckb-backend subscription race); calling Settle again to ForceClose", c.ch.Phase(), wait)
		if err := c.ch.Settle(context.TODO(), false); err != nil {
			panic(err)
		}
	}
	log.Printf("==> Force settle took %v", time.Since(start).Round(time.Millisecond))
	c.ch.Close()
}
