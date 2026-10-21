#!/usr/bin/env bash
# Lab 05: start / stop / inspect the payment-channel-ckb devnet.
# start does exactly what `make dev` does, but in the background (detached tmux),
# and waits until the contracts are deployed and SUDT has been issued.
set -uo pipefail
PERUN="${PERUN:-/root/perun-examples}"
DEV="$PERUN/payment-channel-ckb/devnet"
RPC=http://127.0.0.1:8114
LOGS="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/logs"
mkdir -p "$LOGS"

rpc() {
  curl -s -m 3 -X POST -H 'content-type: application/json' \
    -d "{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"$1\",\"params\":${2:-[]}}" $RPC
}
tip() { rpc get_tip_block_number | jq -r '.result // empty' 2>/dev/null; }

status() {
  if tmux has-session -t devnet 2>/dev/null; then echo "tmux session devnet: running"; else echo "tmux session devnet: not running"; fi
  t=$(tip); [ -n "$t" ] && echo "RPC $RPC: block $((t))" || echo "RPC $RPC: not responding"
  ls "$DEV"/contracts/migrations/dev/*.json >/dev/null 2>&1 && echo "Contracts: deployed" || echo "Contracts: not deployed"
  [ -f "$DEV/accounts/alice.pk" ] && echo "Accounts: $(ls "$DEV/accounts" | tr '\n' ' ')"
}

case "${1:-status}" in
start)
  for t in ckb ckb-cli tmux tmuxp expect python; do
    command -v $t >/dev/null || { echo "Missing $t. Run: ./run.sh lab05-tools"; exit 1; }
  done
  if tmux has-session -t devnet 2>/dev/null; then
    echo "Devnet is already running. For a fresh one: ./run.sh lab05-devnet stop, then start again."; status; exit 0
  fi
  echo ">>> setup-devnet.sh: build contracts, create accounts, ckb init (log: $LOGS/devnet-setup.log)"
  echo "    Every start DELETES the old accounts/, data/, specs/ and creates a fresh devnet."
  cd "$DEV" || exit 1
  # tmuxp runs before_script (setup-devnet.sh), then opens the session in the background (-d).
  if ! tmuxp load -d ./devnet-session.yaml 2>&1 | tee "$LOGS/devnet-setup.log"; then
    echo "setup-devnet.sh failed, see $LOGS/devnet-setup.log"; exit 1
  fi
  echo ">>> Waiting for ckb to run, fund Alice/Bob, deploy contracts, issue SUDT (usually 40–90 seconds)"
  for i in $(seq 180); do
    n=$(ls contracts/migrations/dev/*.json contracts/migrations_vc/dev/*.json 2>/dev/null | wc -l)
    if [ "$n" -ge 2 ] && tmux capture-pane -p -t devnet:1 2>/dev/null | grep -q "listing SUDT account balances"; then
      sleep 8; break
    fi
    printf '.'; sleep 2
  done
  echo
  # The session's own `sudt_helper.sh fund` runs only 1 s after deploying, often before the deploy
  # transactions are committed, and then silently issues nothing. Retry until Alice holds SUDT.
  sudt_alice() { ./sudt_helper.sh balances 2>/dev/null | awk '/total_amount/ {gsub(/"/, "", $2); print $2; exit}'; }
  for attempt in 1 2 3; do
    [ "$(sudt_alice)" != "0" ] && [ -n "$(sudt_alice)" ] && break
    echo ">>> Alice has no SUDT yet, issuing SUDT (attempt $attempt)"
    ./sudt_helper.sh fund >>"$LOGS/devnet-setup.log" 2>&1
    for _ in $(seq 15); do sleep 2; [ "$(sudt_alice)" != "0" ] && break; done
  done
  echo ">>> SUDT balance of Alice: $(sudt_alice)"
  echo ">>> Deploy window (tmux devnet:1):"
  tmux capture-pane -p -t devnet:1 | grep -v '^$' | tail -15
  status
  echo
  echo "Watch live: tmux attach -t devnet   (detach: Ctrl+b then d)"
  echo "Next step: ./run.sh lab05"
  ;;
stop)
  tmux kill-session -t devnet 2>/dev/null && echo "Devnet stopped." || echo "Devnet is not running."
  pkill -x ckb 2>/dev/null || true
  ;;
status) status ;;
*) echo "Usage: $0 start|stop|status"; exit 2 ;;
esac
