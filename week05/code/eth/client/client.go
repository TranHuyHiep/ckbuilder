// Based on perun-examples/payment-channel/client/client.go (Apache-2.0, PolyCrypt GmbH).
// Changes for the week 5 labs:
//   - OpenChannel takes peerAmount and challengeDuration and returns an error instead of panicking.
//   - Policy: the receiving side's policy (maximum funding, minimum payment).
//   - Buggy: re-enables the original code's missing `return` after Reject (Lab 03).

package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Policy is what the client enforces on incoming proposals and updates.
type Policy struct {
	// MaxPeerFunding is the most wei we agree to deposit when invited into a channel.
	// The original code is equivalent to MaxPeerFunding = 0.
	MaxPeerFunding *big.Int
	// MinPayment is the smallest payment we accept (anti-spam).
	// nil means no limit.
	MinPayment *big.Int
}

// DefaultPolicy matches the original code: deposit nothing, accept any payment.
func DefaultPolicy() Policy {
	return Policy{MaxPeerFunding: big.NewInt(0)}
}

// PaymentClient is a payment channel client.
type PaymentClient struct {
	Name        string                              // Name shown in logs (Alice, Bob).
	Policy      Policy                              // Policy for incoming proposals and updates.
	Buggy       bool                                // true = keep the original missing-`return` bug.
	perunClient *client.Client                      // The core Perun client.
	account     map[wallet.BackendID]wallet.Address // The account we use for on-chain and off-chain transactions.
	waddress    map[wallet.BackendID]wire.Address
	currency    channel.Asset        // The currency we expect to get paid in.
	channels    chan *PaymentChannel // Accepted payment channels.
}

// SetupPaymentClient creates a new payment client.
func SetupPaymentClient(
	name string,
	bus wire.Bus,
	w *swallet.Wallet,
	acc common.Address,
	eaddress *ethwallet.Address,
	nodeURL string,
	chainID uint64,
	adjudicator common.Address,
	assetaddr ethwallet.Address,
	wireAddr wire.Address,
) (*PaymentClient, error) {
	// Create Ethereum client and contract backend.
	cb, err := CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, fmt.Errorf("creating contract backend: %w", err)
	}

	// Validate contracts.
	err = ethchannel.ValidateAdjudicator(context.TODO(), cb, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}
	err = ethchannel.ValidateAssetHolderETH(context.TODO(), cb, common.Address(assetaddr), adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}

	// Setup funder.
	funder := ethchannel.NewFunder(cb)
	dep := ethchannel.NewETHDepositor(50000)
	ethAcc := accounts.Account{Address: acc}
	asset := ethchannel.NewAsset(big.NewInt(int64(chainID)), common.Address(assetaddr))
	funder.RegisterAsset(*asset, dep, ethAcc)

	// Setup adjudicator.
	adj := ethchannel.NewAdjudicator(cb, adjudicator, acc, ethAcc, 1000000)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	wireAddrs := map[wallet.BackendID]wire.Address{ethwallet.BackendID: wireAddr}
	perunClient, err := client.New(wireAddrs, bus, funder, adj, map[wallet.BackendID]wallet.Wallet{ethwallet.BackendID: w}, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	eAddrs := map[wallet.BackendID]wallet.Address{ethwallet.BackendID: eaddress}

	// Create client and start request handler.
	c := &PaymentClient{
		Name:        name,
		Policy:      DefaultPolicy(),
		perunClient: perunClient,
		account:     eAddrs,
		waddress:    wireAddrs,
		currency:    asset,
		channels:    make(chan *PaymentChannel, 1),
	}
	go perunClient.Handle(c, c)

	return c, nil
}

// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(
	peer map[wallet.BackendID]wire.Address,
	amount, peerAmount float64, // ETH we deposit and ETH the peer must deposit.
	challengeDuration uint64, // On-chain challenge duration in seconds.
) (*PaymentChannel, error) {
	// The proposer has always index 0.
	participants := []map[wallet.BackendID]wire.Address{c.waddress, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, []wallet.BackendID{ethwallet.BackendID}, c.currency)
	initAlloc.SetAssetBalances(c.currency, []channel.Bal{
		EthToWei(big.NewFloat(amount)),     // Our initial balance.
		EthToWei(big.NewFloat(peerAmount)), // Peer's initial balance.
	})

	// Prepare the channel proposal by defining the channel parameters.
	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account,
		initAlloc,
		participants,
	)
	if err != nil {
		return nil, fmt.Errorf("creating proposal: %w", err)
	}

	// Send the proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		return nil, err
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	return newPaymentChannel(ch, c.currency), nil
}

// startWatching starts the dispute watcher for the specified channel.
func (c *PaymentClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			log.Printf("[%s] watcher returned with error: %v", c.Name, err)
		}
	}()
}

// AcceptedChannel returns the next accepted channel.
func (c *PaymentClient) AcceptedChannel() *PaymentChannel {
	return <-c.channels
}

// AcceptedChannelTimeout is like AcceptedChannel but returns false on timeout,
// for when the proposal may be rejected.
func (c *PaymentClient) AcceptedChannelTimeout(d time.Duration) (*PaymentChannel, bool) {
	select {
	case ch := <-c.channels:
		return ch, true
	case <-time.After(d):
		return nil, false
	}
}

// Shutdown gracefully shuts down the client.
func (c *PaymentClient) Shutdown() {
	c.perunClient.Close()
}
