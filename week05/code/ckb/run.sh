#!/usr/bin/env bash
# Lab 05: runs one scenario on the CKB devnet. Usage: run.sh default|sudt|inspect|dispute
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOGS="$HERE/../logs"
mkdir -p "$LOGS"
sc=${1:-default}
LOG="$LOGS/lab05-$sc-$(date +%H%M%S).log"

st=$(bash "$HERE/devnet.sh" status)   # capture first: grep -q on a pipe + pipefail = SIGPIPE false negatives
grep -q "Contracts: deployed" <<<"$st" \
  || { echo "Devnet is not ready. Run: ./run.sh lab05-devnet start"; exit 1; }

cd "$HERE"   # util.go reads ./devnet/... (symlink to perun-examples/payment-channel-ckb/devnet)
echo ">>> go run . -scenario $sc   (log: $LOG)"

if [ "$sc" = inspect ]; then
  go run . -scenario inspect -pause 30s > >(tee "$LOG") 2>&1 &
  pid=$!
  for _ in $(seq 300); do grep -q CHANNEL_OPEN "$LOG" 2>/dev/null && break; sleep 1; done
  echo; echo ">>> Channel is open, inspecting cells on-chain:"
  bash "$HERE/cells.sh" | tee -a "$LOG"
  wait $pid
  echo; echo ">>> After settling, inspecting again (the channel cell and funds cells should be gone):"
  sleep 3
  bash "$HERE/cells.sh" | tee -a "$LOG"
else
  go run . -scenario "$sc" 2>&1 | tee "$LOG"
fi
