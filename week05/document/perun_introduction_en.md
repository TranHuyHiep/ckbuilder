# Perun Framework — Introduction (Study Notes, English)

> **Source:** [Perun Framework: Introduction](https://hyperledger-labs.github.io/perun-doc/introduction.html)
> **Vietnamese version:** [perun_introduction_vi.md](./perun_introduction_vi.md)

---

## 1. What is Perun?

Perun is an **open-source project that increases blockchain transaction throughput**.

- It uses **only a few on-chain transactions** to set things up, then moves most activity into an **off-chain peer-to-peer network**.
- After participants open some basic channels, **anyone can transact with anyone else through virtual channels**, with no extra on-chain setup.
- **Goal:** make blockchains ready for mass adoption by fixing three problems:
  - high transaction **fees**,
  - high **latency** (slow confirmation),
  - low **throughput** (few transactions per second).

### Origin
- Started as an independent open-source collaboration between:
  - **Robert Bosch GmbH**, through its "Economy of Things" project, and
  - the **Perun team at the Technical University of Darmstadt (TUDa)**.
- Later joined the **Hyperledger Foundation** as a **Hyperledger Labs** project.

---

## 2. The Perun Protocol

The framework is built on the **Perun protocol**. It is a family of cryptographic protocols designed and **formally analyzed** by researchers at TUDa and the University of Warsaw.

### How it works
1. **Open a channel.** Users lock coins in a smart contract on the blockchain.
2. **Update off-chain.** Users exchange signed state updates directly with each other. No blockchain transaction is needed for each update.
3. **Enforce on-chain if needed.** The blockchain guarantees that the latest off-chain state can always be enforced on-chain, for example when a party cheats or goes offline.
4. **Close/settle.** The final balances are written back to the blockchain.

Both **payment channels** (only balances) and **state channels** (arbitrary app state and smart contract logic) are supported.

### Key features compared with other channel systems such as the Lightning Network

| Feature | Meaning |
|---|---|
| **State-channel virtualization** | Existing channels can be combined into **virtual channels**. Two users without a direct channel can still transact through intermediaries. Virtual channels open and close **fully off-chain**. *The doc says this is not yet implemented.* |
| **Blockchain-agnostic** | A modular design lets Perun plug into **any blockchain** or even a traditional ledger. |
| **Interoperability** | Agnostic design plus virtualization enables **cross-chain** transactions and smart contract execution. Logging, messaging and persistence can be customized. |
| **High security** | The protocol specifications are **mathematically proven** secure with modern research methods. |

### Use cases
Finance/FinTech, mobility, energy, e-commerce, telecommunications, and any scenario that needs **direct micro-transactions**.

### Publications
- IEEE: https://ieeexplore.ieee.org/document/8835315
- ACM (CCS 2018): https://dl.acm.org/doi/10.1145/3243734.3243856
- Springer, *Multi-party Virtual State Channels*: https://www.springerprofessional.de/en/multi-party-virtual-state-channels/16720256

---

## 3. Framework Components

The Perun Framework has **three main software components**:

```
             ┌───────────────────────────┐
  user ───▶  │        perun-node         │  multi-user node (gRPC API)
             │  keys, IDs, networking,   │  ◀── perun-nodecli (CLI client)
             │  open / transact / settle │
             └────────────┬──────────────┘
                          │ uses
             ┌────────────▼──────────────┐
             │         go-perun          │  SDK: proposal protocol,
             │  (blockchain-agnostic)    │  state machine + persistence,
             └────────────┬──────────────┘  watcher, chain connectors
                          │ talks to
             ┌────────────▼──────────────┐
             │   perun-eth-contracts     │  on-chain smart contracts
             │     (Ethereum)            │
             └───────────────────────────┘
```

1. **perun-eth-contracts:** the Ethereum smart contracts needed to run the Perun protocol on-chain.
2. **go-perun:** a Go SDK that implements the core protocol:
   - the **state-channel proposal protocol**, which handles how channels are proposed and opened,
   - a **state machine with persistence**, which tracks channel state and saves it to survive restarts,
   - a **watcher**, which monitors the chain and reacts if a counterparty tries to register an old state.
   - It ships an Ethereum connector but is designed to be **blockchain-agnostic**.
3. **perun-node:** a **multi-user node** built on go-perun. It lets users:
   - manage **keys and identities**,
   - handle **off-chain networking** with peers,
   - **open, transact in, and settle** channels.

### Current status
| Component | Status |
|---|---|
| go-perun SDK | Fully generalized **state channels on Ethereum**, with persistence |
| perun-node | **Two-party payment channels**, local file-based ID provider, **gRPC** interface |
| perun-nodecli | Reference client for the perun-node gRPC API, interactive CLI with auto-completion |

### Roadmap
1. Virtual channels
2. SSI (Self-Sovereign Identity) integration with Hyperledger Aries
3. More blockchain backends
4. Cross-chain channels

---

## 4. Key Takeaways

- Perun is a **Layer-2 channel framework**: lock funds on-chain once, transact off-chain many times, settle on-chain once.
- Security comes from the fact that **off-chain states are always enforceable on-chain**.
- Its distinguishing ideas are **virtual channels** and **blockchain-agnosticism**.
- The software stack is **contracts → go-perun SDK → perun-node**, from the chain up to the user.

## 5. Glossary

| Term | Definition |
|---|---|
| Off-chain | Activity that happens between users without writing to the blockchain |
| Payment channel | A channel that only tracks balances between participants |
| State channel | A channel that tracks arbitrary application state, not just balances |
| Virtual channel | A channel built on top of existing channels through an intermediary, opened and closed without on-chain transactions |
| Watcher | A component that watches the chain and disputes outdated states on the user's behalf |
| Settle | Finalize a channel and pay out the final balances on-chain |

## 6. Further Reading
- Next page in the docs: [Concepts: State Channels](https://hyperledger-labs.github.io/perun-doc/concepts/state-channel.html)
- [Protocols](https://hyperledger-labs.github.io/perun-doc/concepts/protocols.html)
- GitHub: [go-perun](https://github.com/hyperledger-labs/go-perun), [perun-node](https://github.com/hyperledger-labs/perun-node), [perun-eth-contracts](https://github.com/hyperledger-labs/perun-eth-contracts)
