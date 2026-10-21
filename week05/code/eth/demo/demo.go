// Package demo bundles the environment setup of perun-examples/payment-channel
// (main.go + util.go, Apache-2.0, PolyCrypt GmbH) so the labs can reuse it.
package demo

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethchannel "github.com/perun-network/perun-eth-backend/channel"
	ethwallet "github.com/perun-network/perun-eth-backend/wallet"
	swallet "github.com/perun-network/perun-eth-backend/wallet/simple"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net"
	p2p "perun.network/go-perun/wire/net/libp2p"
	perunio "perun.network/go-perun/wire/perunio/serializer"

	"week05/eth/client"
)

const (
	chainID = 1337

	// Private keys, matching hardhat.config.js.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

// ChainURL reads the CHAIN_URL environment variable, default ws://127.0.0.1:8545.
func ChainURL() string {
	if u := os.Getenv("CHAIN_URL"); u != "" {
		return u
	}
	return "ws://127.0.0.1:8545"
}

// Env holds the two ready clients and a wallet balance reader.
type Env struct {
	Alice, Bob *client.PaymentClient
	eth        *ethclient.Client
}

// Setup deploys the contracts, sets up the wire bus and creates clients Alice and Bob.
func Setup() *Env {
	url := ChainURL()

	log.Println("Deploying contracts.")
	adjudicator, assetHolder := deployContracts(url, keyDeployer)
	asset := *ethwallet.AsWalletAddr(assetHolder)

	aliceWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	aliceBus, aliceDialer := setupBusWire(aliceWireAcc)
	bobWireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	bobBus, _ := setupBusWire(bobWireAcc)
	aliceDialer.Register(map[wallet.BackendID]wire.Address{ethwallet.BackendID: bobWireAcc.Address()}, bobWireAcc.ID().String())

	log.Println("Setting up clients.")
	alice := setupPaymentClient("Alice", aliceBus, url, adjudicator, asset, keyAlice, aliceWireAcc.Address())
	bob := setupPaymentClient("Bob", bobBus, url, adjudicator, asset, keyBob, bobWireAcc.Address())

	c, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return &Env{Alice: alice, Bob: bob, eth: c}
}

// Wallets returns the on-chain wallet balances [Alice, Bob] in ETH.
func (e *Env) Wallets() []*big.Float {
	res := []*big.Float{}
	for _, a := range []common.Address{e.Alice.WalletAddress(), e.Bob.WalletAddress()} {
		bal, err := e.eth.BalanceAt(context.TODO(), a, nil)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, client.WeiToEth(bal))
	}
	return res
}

// LogWallets prints the on-chain wallet balances with a label.
func (e *Env) LogWallets(label string) {
	w := e.Wallets()
	log.Printf("Wallets %-7s Alice=%s ETH  Bob=%s ETH", label, Fmt(w[0]), Fmt(w[1]))
}

// Shutdown closes both clients.
func (e *Env) Shutdown() {
	e.Alice.Shutdown()
	e.Bob.Shutdown()
}

// Fmt formats ETH with 6 decimals.
func Fmt(f *big.Float) string {
	return f.Text('f', 6)
}

// Must exits the program if err is not nil.
func Must(err error, what string) {
	if err != nil {
		log.Fatalf("%s: %v", what, err)
	}
}

// Header prints a title line to make logs easier to read.
func Header(s string) {
	fmt.Printf("\n========== %s ==========\n", s)
}

// deployContracts deploys the Perun smart contracts on the specified ledger.
func deployContracts(nodeURL string, privateKey string) (adj, ah common.Address) {
	k, err := crypto.HexToECDSA(privateKey)
	Must(err, "parsing key")
	w := swallet.NewWallet(k)
	cb, err := client.CreateContractBackend(nodeURL, chainID, w)
	Must(err, "creating contract backend (is Hardhat running?)")
	acc := accounts.Account{Address: crypto.PubkeyToAddress(k.PublicKey)}

	adj, err = ethchannel.DeployAdjudicator(context.TODO(), cb, acc)
	Must(err, "deploying adjudicator")
	ah, err = ethchannel.DeployETHAssetholder(context.TODO(), cb, adj, acc)
	Must(err, "deploying asset holder")
	return adj, ah
}

// setupPaymentClient sets up a new client with the given parameters.
func setupPaymentClient(
	name string,
	bus wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	asset ethwallet.Address,
	privateKey string,
	wireAddr wire.Address,
) *client.PaymentClient {
	k, err := crypto.HexToECDSA(privateKey)
	Must(err, "parsing key")
	w := swallet.NewWallet(k)
	acc := crypto.PubkeyToAddress(k.PublicKey)
	eaddr := ethwallet.AsWalletAddr(acc)

	c, err := client.SetupPaymentClient(name, bus, w, acc, eaddr, nodeURL, chainID, adjudicator, asset, wireAddr)
	Must(err, "setting up client "+name)
	return c
}

// setupBusWire sets up a wire.Bus for the given wire.Account.
func setupBusWire(acc *p2p.Account) (wire.Bus, *p2p.Dialer) {
	id := make(map[wallet.BackendID]wire.Account)
	id[ethwallet.BackendID] = acc

	listener := p2p.NewP2PListener(acc)
	dialer := p2p.NewP2PDialer(acc)

	bus := net.NewBus(id, dialer, perunio.Serializer())

	go bus.Listen(listener)
	return bus, dialer
}
