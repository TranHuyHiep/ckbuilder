#!/usr/bin/env python3
"""Reads an `npx hardhat node` log and summarizes the on-chain transactions.

Usage: python3 txsummary.py <hardhat.log>
"""
import re
import sys

NAMES = {
    "0xe84d227431dffcf14fb8fa39818dfd4e864aeb13": "Deployer",
    "0x56fd289cee714a5e471c418436efa63e780d7a87": "Alice",
    "0x6536425be95a6661f6c6f68d709b6be152785df6": "Bob",
}
ANSI = re.compile(r"\x1b\[[0-9;]*m")


def parse(path):
    txs, cur = [], None
    for raw in open(path, encoding="utf-8", errors="replace"):
        line = ANSI.sub("", raw).rstrip()
        if line.startswith("eth_sendRawTransaction"):
            cur = {}
            txs.append(cur)
            continue
        if cur is None:
            continue
        if not line.startswith("  "):
            cur = None
            continue
        key, _, val = line.strip().partition(":")
        cur[key.strip()] = val.strip()
    return [t for t in txs if "Transaction" in t]


def main():
    txs = parse(sys.argv[1])
    deployed = [t["Contract address"].lower() for t in txs if "Contract deployment" in t]
    role = {}
    if len(deployed) >= 1:
        role[deployed[0]] = "Adjudicator"
    if len(deployed) >= 2:
        role[deployed[1]] = "AssetHolder"

    counts = {}
    print("\nOn-chain transactions (from the Hardhat log):")
    print("| # | Block | Sender | Target | Value | Gas | Kind |")
    print("|---|---|---|---|---|---|---|")
    for i, t in enumerate(txs, 1):
        frm = NAMES.get(t.get("From", "").lower(), t.get("From", "?"))
        to = t.get("To", "").lower()
        value = t.get("Value", "0 ETH").replace(". ETH", " ETH")
        block = re.sub(r":.*", "", next((k for k in t if k.startswith("Block #")), "Block ?"))
        if "Contract deployment" in t:
            dest = role.get(t["Contract address"].lower(), "?")
            kind = "deploy " + dest
        else:
            dest = role.get(to, to[:10])
            if dest == "AssetHolder":
                kind = "deposit" if not value.startswith("0 ") else "withdraw"
            elif dest == "Adjudicator":
                kind = "register / conclude"
            else:
                kind = "other"
        counts[kind] = counts.get(kind, 0) + 1
        print(f"| {i} | {block.replace('Block ', '')} | {frm} | {dest} | {value} | {t.get('Gas used', '?').split(' of')[0]} | {kind} |")
    print("\nTotal: %d transactions. " % len(txs) + ", ".join(f"{k}: {v}" for k, v in counts.items()))


if __name__ == "__main__":
    main()
