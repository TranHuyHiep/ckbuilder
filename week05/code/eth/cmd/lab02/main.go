// Lab 02: runs the original Alice ↔ Bob scenario and prints the channel lifecycle
// table (in-channel balances + version after each step), with the time per step.
package main

import (
	"fmt"
	"log"
	"time"

	"week05/eth/client"
	"week05/eth/demo"
)

func row(step int, action string, ch *client.PaymentChannel, d time.Duration) {
	b := ch.Balances()
	fmt.Printf("| %d | %-28s | %5s | %5s | %d | %-5v | %8v |\n",
		step, action, b[0].Text('f', 2), b[1].Text('f', 2), ch.Version(), ch.IsFinal(), d.Round(time.Millisecond))
}

func main() {
	env := demo.Setup()
	defer env.Shutdown()
	env.LogWallets("before")

	demo.Header("Open channel")
	t := time.Now()
	chAlice, err := env.Alice.OpenChannel(env.Bob.WireAddress(), 5, 0, 10)
	demo.Must(err, "OpenChannel")
	chBob := env.Bob.AcceptedChannel()
	dOpen := time.Since(t)

	fmt.Println("\n| Step | Action | Alice (channel) | Bob (channel) | Version | IsFinal | Time |")
	fmt.Println("|---|---|---|---|---|---|---|")
	row(0, "OpenChannel(bob, 5)", chAlice, dOpen)

	t = time.Now()
	demo.Must(chAlice.SendPayment(3), "Alice pay 3")
	row(1, "chAlice.SendPayment(3)", chAlice, time.Since(t))

	t = time.Now()
	demo.Must(chBob.SendPayment(1), "Bob pay 1")
	row(2, "chBob.SendPayment(1)", chAlice, time.Since(t))

	t = time.Now()
	demo.Must(chAlice.SendPayment(1), "Alice pay 1")
	row(3, "chAlice.SendPayment(1)", chAlice, time.Since(t))

	t = time.Now()
	demo.Must(chAlice.Finalize(), "Finalize")
	row(4, "Finalize (IsFinal = true)", chAlice, time.Since(t))
	fmt.Println()

	demo.Header("Settle")
	t = time.Now()
	demo.Must(chAlice.Settle(), "Alice settle") // Conclude and withdraw.
	demo.Must(chBob.Settle(), "Bob settle")     // Withdraw.
	log.Printf("Settle (conclude + 2 withdraw) took %v", time.Since(t).Round(time.Millisecond))

	env.LogWallets("after")
}
