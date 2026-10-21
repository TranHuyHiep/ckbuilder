// Lab 03: attacking the channel, the missing-`return` bug, two-sided funding, anti-spam rule.
// Run: go run ./cmd/lab03 -scenario <name>. List scenarios with -h.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"week05/eth/client"
	"week05/eth/demo"
)

var scenarios = []struct{ name, desc string }{
	{"mint", "Alice adds 100 ETH to herself (debiting no one)"},
	{"steal-buggy", "Alice takes 2 ETH from Bob, Bob runs the ORIGINAL buggy code"},
	{"steal", "Alice takes 2 ETH from Bob, Bob runs the FIXED code"},
	{"fund52-buggy", "Alice asks Bob to deposit 2 ETH, original policy (Bob deposits 0), ORIGINAL buggy code"},
	{"fund52-default", "Alice asks Bob to deposit 2 ETH, original policy (Bob deposits 0), FIXED code"},
	{"fund52", "Deposit 5/2, new policy: Bob deposits at most 3 ETH, full scenario"},
	{"fund54", "Deposit 5/4, new policy: Bob deposits at most 3 ETH"},
	{"spam", "Bob rejects payments < 0.001 ETH; Alice pays 0.0001 then 0.01"},
	{"afterfinal", "Send a payment after the state is final"},
	{"micro", "100 payments of 0.01 ETH, timed"},
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/lab03 -scenario <name>\n\nScenarios:")
	for _, s := range scenarios {
		fmt.Fprintf(os.Stderr, "  %-15s %s\n", s.name, s.desc)
	}
}

func eth(x float64) *big.Int { return client.EthToWei(big.NewFloat(x)) }

func logCh(label string, ch *client.PaymentChannel) {
	b := ch.Balances()
	log.Printf("In channel %-24s Alice=%s Bob=%s version=%d final=%v",
		label, b[0].Text('f', 4), b[1].Text('f', 4), ch.Version(), ch.IsFinal())
}

// open opens a channel and waits for Bob to accept. Returns nil if rejected.
func open(env *demo.Env, amount, peerAmount float64) (*client.PaymentChannel, *client.PaymentChannel) {
	log.Printf("Alice proposes a channel: Alice deposits %v ETH, Bob deposits %v ETH", amount, peerAmount)
	chAlice, err := env.Alice.OpenChannel(env.Bob.WireAddress(), amount, peerAmount, 10)
	if err != nil {
		log.Printf("==> Alice got an error opening the channel: %v", err)
		return nil, nil
	}
	chBob, ok := env.Bob.AcceptedChannelTimeout(10 * time.Second)
	if !ok {
		log.Fatal("Bob did not receive the channel")
	}
	logCh("after open", chAlice)
	return chAlice, chBob
}

func settle(env *demo.Env, chAlice, chBob *client.PaymentChannel) {
	demo.Must(chAlice.Settle(), "Alice settle")
	demo.Must(chBob.Settle(), "Bob settle")
	env.LogWallets("after")
}

func main() {
	name := flag.String("scenario", "", "scenario name")
	flag.Usage = usage
	flag.Parse()
	desc := ""
	for _, s := range scenarios {
		if s.name == *name {
			desc = s.desc
		}
	}
	if desc == "" {
		usage()
		os.Exit(2)
	}
	demo.Header(*name + ": " + desc)

	env := demo.Setup()
	defer env.Shutdown()
	env.LogWallets("before")

	switch *name {
	case "mint":
		chAlice, chBob := open(env, 5, 0)
		err := chAlice.MintForMe(100)
		log.Printf("==> MintForMe(100) returned: %v", err)
		logCh("after MintForMe", chAlice)
		settle(env, chAlice, chBob)

	case "steal-buggy", "steal":
		env.Bob.Buggy = *name == "steal-buggy"
		chAlice, chBob := open(env, 5, 0)
		demo.Must(chAlice.SendPayment(3), "Alice pay 3")
		logCh("after SendPayment(3)", chAlice)
		log.Println("Alice tries Steal(2)...")
		err := chAlice.Steal(2)
		log.Printf("==> Steal(2) returned: %v", err)
		logCh("after Steal", chAlice)
		settle(env, chAlice, chBob)
		log.Println("==> Program ran to the end, no panic.")

	case "fund52-buggy", "fund52-default":
		env.Bob.Buggy = *name == "fund52-buggy"
		chAlice, chBob := open(env, 5, 2)
		if chAlice != nil {
			settle(env, chAlice, chBob)
		}
		time.Sleep(time.Second) // Give Bob's goroutine time to run (and panic, if buggy).
		log.Println("==> Program ran to the end, no panic.")

	case "fund52", "fund54":
		env.Bob.Policy.MaxPeerFunding = eth(3)
		peer := 2.0
		if *name == "fund54" {
			peer = 4
		}
		chAlice, chBob := open(env, 5, peer)
		if chAlice == nil {
			break
		}
		demo.Must(chAlice.SendPayment(3), "Alice pay 3")
		demo.Must(chBob.SendPayment(1), "Bob pay 1")
		demo.Must(chAlice.SendPayment(1), "Alice pay 1")
		logCh("after 3 payments", chAlice)
		settle(env, chAlice, chBob)

	case "spam":
		env.Bob.Policy.MinPayment = eth(0.001)
		chAlice, chBob := open(env, 5, 0)
		log.Printf("==> SendPayment(0.0001) returned: %v", chAlice.SendPayment(0.0001))
		log.Printf("==> SendPayment(0.01)   returned: %v", chAlice.SendPayment(0.01))
		logCh("after 2 attempts", chAlice)
		settle(env, chAlice, chBob)

	case "afterfinal":
		chAlice, chBob := open(env, 5, 0)
		demo.Must(chAlice.SendPayment(3), "Alice pay 3")
		demo.Must(chAlice.Finalize(), "Finalize")
		logCh("after Finalize", chAlice)
		log.Printf("==> SendPayment(1) after final returned: %v", chAlice.SendPayment(1))
		settle(env, chAlice, chBob)

	case "micro":
		chAlice, chBob := open(env, 5, 0)
		const n = 100
		start := time.Now()
		for i := 0; i < n; i++ {
			demo.Must(chAlice.SendPayment(0.01), "micro payment")
		}
		d := time.Since(start)
		logCh("after 100 payments", chAlice)
		log.Printf("==> %d off-chain payments took %v (%.2f ms each)", n, d.Round(time.Millisecond), float64(d.Microseconds())/1000/n)
		log.Printf("==> If each were one 5-second on-chain block: at least %v", time.Duration(n)*5*time.Second)
		settle(env, chAlice, chBob)
	}
}
