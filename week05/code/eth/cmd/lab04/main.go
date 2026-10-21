// Lab 04: compares the fast (cooperative) settle path with the dispute (uncooperative) path.
// Run: go run ./cmd/lab04 -mode fast|dispute [-challenge 10]
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"week05/eth/demo"
)

func main() {
	mode := flag.String("mode", "dispute", "fast = finalize then conclude; dispute = register, wait for challenge, conclude")
	challenge := flag.Uint64("challenge", 10, "challenge duration (seconds)")
	flag.Parse()
	if *mode != "fast" && *mode != "dispute" {
		flag.Usage()
		os.Exit(2)
	}
	demo.Header("mode=" + *mode)

	env := demo.Setup()
	defer env.Shutdown()
	env.LogWallets("before")

	log.Printf("Opening channel with challenge duration = %d seconds", *challenge)
	chAlice, err := env.Alice.OpenChannel(env.Bob.WireAddress(), 5, 0, *challenge)
	demo.Must(err, "OpenChannel")
	chBob := env.Bob.AcceptedChannel()

	demo.Must(chAlice.SendPayment(3), "Alice pay 3")
	demo.Must(chBob.SendPayment(1), "Bob pay 1")
	demo.Must(chAlice.SendPayment(1), "Alice pay 1")
	b := chAlice.Balances()
	log.Printf("In channel: Alice=%s Bob=%s version=%d final=%v", demo.Fmt(b[0]), demo.Fmt(b[1]), chAlice.Version(), chAlice.IsFinal())

	start := time.Now()
	if *mode == "fast" {
		log.Println("Alice: Finalize + Settle (fast path)...")
		demo.Must(chAlice.Settle(), "Alice settle")
		log.Printf("==> Alice settle took %v", time.Since(start).Round(time.Millisecond))
		demo.Must(chBob.Settle(), "Bob settle")
	} else {
		log.Println("Alice: ForceSettle, NO finalize (simulating an uncooperative Bob)...")
		d, err := chAlice.ForceSettle()
		demo.Must(err, "Alice force settle")
		log.Printf("==> Alice force settle took %v (challenge duration %ds, block time 5s)", d.Round(time.Millisecond), *challenge)
		demo.Must(chBob.WithdrawOnly(), "Bob withdraw")
	}
	log.Printf("==> Total settle time for both: %v", time.Since(start).Round(time.Millisecond))
	time.Sleep(500 * time.Millisecond) // Let the last event get printed.
	env.LogWallets("after")
}
