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

// Lab 05: based on perun-examples/payment-channel-ckb/main.go, with scenarios
// selected by -scenario:
//
//	default  original scenario: open a 100/100 CKB channel, Alice pays 10, Bob pays 5, settle
//	sudt     like default but the channel holds both CKB and SUDT, Alice also pays 10 SUDT
//	inspect  like default but pauses -pause before settling so you can inspect cells on-chain
//	dispute  like default but Alice settles WITHOUT finalizing (dispute path)
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
	"github.com/perun-network/perun-libp2p-wire/p2p"
	"perun.network/go-perun/channel"
	ckbchannel "perun.network/perun-ckb-backend/channel"
	"perun.network/perun-ckb-backend/channel/asset"
	"week05/ckb/client"
)

const (
	rpcNodeURL = "http://localhost:8114"
	Network    = types.NetworkTest // Network to use for the CKB client
)

func newClient(name string, idx int, setup *Setup) (*client.PaymentClient, *p2p.Account) {
	wireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))
	net, err := p2p.NewP2PBus(ckbchannel.CKBBackendID, wireAcc)
	if err != nil {
		log.Fatalf("creating p2p net: %v", err)
	}
	go net.Bus.Listen(net.Listener)

	c, err := client.NewPaymentClient(
		name,
		Network,
		setup.Deployment,
		rpcNodeURL,
		setup.WalletAccs[idx],
		*setup.AccKeys[idx],
		setup.Wallets[idx],
		wireAcc.Address(),
		net,
	)
	if err != nil {
		log.Fatalf("error creating %s's client: %v", name, err)
	}
	return c, wireAcc
}

func main() {
	scenario := flag.String("scenario", "default", "default | sudt | inspect | dispute")
	pause := flag.Duration("pause", 60*time.Second, "pause before settling in the inspect scenario")
	flag.Parse()
	switch *scenario {
	case "default", "sudt", "inspect", "dispute":
	default:
		flag.Usage()
		os.Exit(2)
	}
	fmt.Printf("\n========== Lab 05, scenario %s ==========\n", *scenario)

	// Setup devnet environment
	setup := NewSetup()

	log.Println("Setting up payment channel clients")
	alice, _ := newClient("Alice", 0, setup)
	bob, _ := newClient("Bob", 1, setup)

	fmt.Println("Wallet before | Alice:", alice.GetBalances())
	fmt.Println("Wallet before | Bob:  ", bob.GetBalances())

	// Each side deposits exactly this (the CKB example's HandleProposal requires equal funding).
	funding := map[channel.Asset]float64{setup.CKBAsset: 100.0}
	assets := []asset.Asset{*setup.CKBAsset}
	if *scenario == "sudt" {
		funding[setup.SudtAsset] = 50 // smallest SUDT unit
		assets = append(assets, *setup.SudtAsset)
	}

	log.Println("Opening channel and depositing funds")
	start := time.Now()
	chAlice := alice.OpenChannel(bob.WireAddress(), bob.PeerID(), funding)
	chBob := bob.AcceptedChannel()
	log.Printf("==> Opening the channel (Start + Fund on-chain) took %v", time.Since(start).Round(time.Millisecond))
	printBalances(chAlice, assets)

	log.Println("Alice pays Bob 10 CKB")
	pay := map[channel.Asset]float64{setup.CKBAsset: 10.0}
	if *scenario == "sudt" {
		log.Println("Alice also pays Bob 10 SUDT")
		pay[setup.SudtAsset] = 10
	}
	chAlice.SendPayment(pay)
	printBalances(chAlice, assets)

	log.Println("Bob pays Alice 5 CKB")
	chBob.SendPayment(map[channel.Asset]float64{setup.CKBAsset: 5.0})
	printBalances(chAlice, assets)

	if *scenario == "inspect" {
		// ckb/run.sh uses this line as the signal to start inspecting cells.
		log.Printf("CHANNEL_OPEN pausing %v to inspect cells on-chain", *pause)
		time.Sleep(*pause)
	}

	start = time.Now()
	if *scenario == "dispute" {
		log.Println("Alice ForceSettle, NO finalize (simulating an uncooperative Bob)")
		chAlice.ForceSettle()
	} else {
		log.Println("Alice Settle (finalize + Close)")
		chAlice.Settle()
		log.Printf("==> Settle took %v", time.Since(start).Round(time.Millisecond))
	}

	time.Sleep(3 * time.Second) // Wait for the indexer to update balances.
	fmt.Println("Wallet after  | Alice:", alice.GetBalances())
	fmt.Println("Wallet after  | Bob:  ", bob.GetBalances())

	alice.Shutdown()
	bob.Shutdown()
	log.Println("Clients shutdown, exiting method")
}
