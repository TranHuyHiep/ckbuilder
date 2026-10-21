#!/usr/bin/env bash
# Prints the code to read for Lab 02 and Lab 05, with paths to open in VS Code.
PERUN="${PERUN:-/root/perun-examples}"
GP="$(go env GOMODCACHE)/perun.network/go-perun@v0.15.0"
C="$PERUN/payment-channel-ckb/devnet/contracts"

show() { printf '\n\033[1;33m--- %s\033[0m\n' "$1"; }

case "${1:-}" in
lab02)
  show "Components created in SetupPaymentClient ($PERUN/payment-channel/client/client.go)"
  grep -nE "CreateContractBackend|NewFunder|NewETHDepositor|RegisterAsset|NewAdjudicator|NewWatcher|client.New\(|Handle\(c, c\)" "$PERUN/payment-channel/client/client.go"
  show "The three handler interfaces PaymentClient implements ($GP/client)"
  grep -n "ProposalHandler interface" -A4 "$GP/client/proposal.go"
  grep -n "UpdateHandler interface" -A4 "$GP/client/update.go"
  grep -n "AdjudicatorEventHandler interface" -A3 "$GP/client/adjudicate.go"
  show "Update: the steps of one payment ($GP/client/update.go)"
  grep -nE "^func \(c \*Channel\) (Update|update|UpdateBy)|machine\.Update|Sig|conn\.Send|Accept|Enable|NewUpdate" "$GP/client/update.go" | head -30
  show "Channel phases ($GP/channel/machine.go)"
  grep -n "Phase = iota" -A12 "$GP/channel/machine.go"
  show "Challenge duration in the example"
  grep -n "challengeDuration" "$PERUN/payment-channel/client/client.go"
  ;;
lab05)
  show "Contracts ($C/contracts)"
  ls "$C/contracts"
  show "Molecule types ($C/crates/perun-common/types.mol)"
  grep -nE "^(table|union|struct|vector|array) " "$C/crates/perun-common/types.mol"
  show "ChannelWitness"
  grep -n "union ChannelWitness" -A8 "$C/crates/perun-common/types.mol"
  show "Functions in PCTS ($C/contracts/perun-channel-typescript/src/lib.rs)"
  grep -n "^pub fn" "$C/contracts/perun-channel-typescript/src/lib.rs"
  show "Branches in PCTS main()"
  sed -n '/^pub fn main/,/^}/p' "$C/contracts/perun-channel-typescript/src/lib.rs"
  show "enum PChannelAction ($C/crates/perun-common/src/channels.rs)"
  sed -n '/^pub enum PChannelAction/,/^}/p' "$C/crates/perun-common/src/channels.rs"
  show "Time check on ForceClose"
  grep -n "verify_time_lock_expired" -B3 -A3 "$C/contracts/perun-channel-typescript/src/lib.rs" | head -20
  show "The two lock scripts (read them in full)"
  wc -l "$C/contracts/perun-channel-lockscript/src/lib.rs" "$C/contracts/perun-funds-lockscript/src/lib.rs"
  ;;
*) echo "Usage: $0 lab02|lab05"; exit 2 ;;
esac
