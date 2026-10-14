# CKB Economic Model

> **Curriculum Reference:** [truthixify/learn-ckb-in-45-minutes / 04-economic-model.md](https://github.com/truthixify/learn-ckb-in-45-minutes/blob/main/notes/04-economic-model.md)

## The Problem: State Bloat
- Most chains charge a one-time gas fee for storage → data lives forever for free
- Full nodes must store all state forever with no ongoing compensation
- State grows endlessly (Ethereum's growing pain)

## CKB's Solution: 1 CKByte = 1 Byte

- To store N bytes on-chain, you must lock N CKBytes in the cell's `capacity`
- CKB is **locked** while the cell exists — can't be spent elsewhere
- Destroy the cell → get your CKB back
- Natural incentive to clean up unused state
- Total on-chain state is **bounded** by total CKB supply

## Two Issuance Sources

### Base Issuance
- Fixed supply, halves every 4 years (like Bitcoin)
- Goes entirely to **miners**
- Eventually drops to zero

### Secondary Issuance
- Fixed amount of CKBytes per year, **forever**
- Distributed based on how CKB is used:

| CKB Usage | Secondary Issuance Goes To |
|---|---|
| Locked in cells (storing state) | **Miners** — effectively "state rent" |
| Deposited in Nervos DAO | **Back to depositor** — cancels out dilution |
| Liquid (not locked, not in DAO) | **Treasury** (currently burned) |

## The Nervos DAO — Inflation Shield

- Deposit CKB into the DAO if you don't need to store data
- Receive secondary issuance proportional to deposit → exactly offsets dilution
- Effectively opts you out of the "state rent" tax

## The Analogy: Property Tax
- **Using land (on-chain storage):** you pay rent (diluted by secondary issuance)
- **Just holding currency:** deposit in DAO, stay undiluted
- **Miners:** get paid ongoing security from both issuance sources

## Design Pressure on dApp Developers
- Storing 1 GB on-chain = locking 1 billion CKB — massive capital cost
- Pushes devs to **minimize on-chain data** (store hashes/proofs, keep bulk data off-chain)
- Incentivizes **cleaning up** unused cells to reclaim CKB
- Common pattern: **users provide their own storage** — each user's cells hold their own state (tokens, NFTs), rather than one contract holding everything
- State is naturally **distributed across user-owned cells** (unlike Ethereum's contract-centric model)

## Key Takeaways
- CKB ties storage to token value — using space costs ongoing opportunity cost
- Secondary issuance is a perpetual "state rent" mechanism that pays miners
- Nervos DAO protects non-storage holders from dilution
- The system creates bounded state growth and sustainable miner incentives long-term

---

## Questions

**Q1:** On Ethereum, a dApp could deploy a contract that stores 1GB of data and only pay gas once. On CKB, what would that cost, and what natural pressure does this create on dApp developers?

<details>
<summary>Answer</summary>
1 GB = 1,000,000,000 bytes = 1,000,000,000 CKB locked up. That's a massive amount of capital sitting idle. This pushes dApp developers to: minimize on-chain data (store hashes/proofs, keep bulk data off-chain), clean up unused cells to reclaim CKB, and design so users provide their own storage via their own cells rather than the developer footing the bill for everyone.
</details>

---

**Q2:** If every CKB holder deposited into the Nervos DAO, who would pay the "state rent" for miners?

<details>
<summary>Hint</summary>
Think about whether it's possible for ALL CKB to be in the DAO. What must exist on-chain for the system to work?
</details>

<details>
<summary>Answer</summary>
It's impossible for all CKB to be in the DAO — cells that store state (dApp data, tokens, scripts) have CKB locked in their capacity. That CKB can't be in the DAO simultaneously. So there will always be CKB locked in cells paying implicit "state rent" to miners. The more on-chain state exists, the more rent miners collect from secondary issuance.
</details>

---

**Q3:** Why is CKB's economic model more sustainable long-term for miners compared to Bitcoin's?

<details>
<summary>Answer</summary>
Bitcoin's base issuance (block rewards) drops to zero over time, leaving miners dependent entirely on transaction fees — which may not be enough to secure the network. CKB has secondary issuance that runs **forever** at a fixed amount per year. As long as there's state stored on-chain, miners get paid from secondary issuance as "state rent." This creates a sustainable security budget regardless of transaction volume.
</details>
