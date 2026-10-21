#!/usr/bin/env bash
# Runs the week 5 Perun exercises. For every Ethereum scenario the script:
#   1. stops any old Hardhat holding port 8545,
#   2. starts a fresh Hardhat (balances back to 10/10),
#   3. runs the Go scenario and saves the log to code/logs/,
#   4. summarizes the on-chain transactions from the Hardhat log, then stops Hardhat.
#
# See: ./run.sh help
set -uo pipefail

CODE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PERUN="${PERUN:-/root/perun-examples}"
HH_DIR="$PERUN/payment-channel"   # holds hardhat.config.js and node_modules
LOGS="$CODE/logs"
mkdir -p "$LOGS"

HH_PID=""

say() { printf '\n\033[1;36m>>> %s\033[0m\n' "$*"; }

stop_port_8545() {
  local pids
  pids=$(ss -ltnpH 'sport = :8545' 2>/dev/null | grep -o 'pid=[0-9]*' | cut -d= -f2 | sort -u)
  if [ -n "$pids" ]; then
    say "Port 8545 is in use (pid $pids), stopping it to start a fresh chain"
    kill $pids 2>/dev/null
    sleep 1
  fi
}

start_hardhat() {
  if [ ! -d "$HH_DIR/node_modules/hardhat" ]; then
    say "Hardhat is not installed, running npm install in $HH_DIR"
    (cd "$HH_DIR" && npm install --no-audit --no-fund) || exit 1
  fi
  stop_port_8545
  HH_LOG="$LOGS/hardhat-$RUN_ID.log"
  say "Starting a fresh Hardhat (log: $HH_LOG)"
  (cd "$HH_DIR" && exec setsid npx hardhat node >"$HH_LOG" 2>&1) &
  HH_PID=$!
  for _ in $(seq 60); do
    grep -q "Started HTTP" "$HH_LOG" 2>/dev/null && return 0
    sleep 0.5
  done
  echo "Hardhat failed to start, see $HH_LOG"; exit 1
}

stop_hardhat() {
  [ -z "$HH_PID" ] && return
  kill -- -"$HH_PID" 2>/dev/null || kill "$HH_PID" 2>/dev/null
  wait "$HH_PID" 2>/dev/null
  HH_PID=""
}
trap stop_hardhat EXIT

# eth_run <lab> <log-tag> <go run args...>
eth_run() {
  local lab=$1 tag=$2; shift 2
  RUN_ID="$lab-$tag-$(date +%H%M%S)"
  start_hardhat
  say "go run ./cmd/$lab $*   (log: $LOGS/$RUN_ID.log)"
  (cd "$CODE/eth" && timeout 300 go run "./cmd/$lab" "$@") 2>&1 | grep -v '^pkg/test: using rootSeed' | tee "$LOGS/$RUN_ID.log"
  local rc=${PIPESTATUS[0]}
  LAST_LOG="$LOGS/$RUN_ID.log"
  [ "$rc" -ne 0 ] && say "Program exited with code $rc (see the log above)"
  python3 -I "$CODE/tools/txsummary.py" "$HH_LOG" | tee -a "$LOGS/$RUN_ID.log"
  stop_hardhat
}

LAB03_ALL="mint steal-buggy steal fund52-buggy fund52-default fund52 fund54 spam afterfinal micro"

usage() {
  cat <<EOF
Usage: ./run.sh <command> [args]

Lab 01
  check                   Check tools and disk, test-build the modules

Lab 02
  lab02                   Run the demo, print the channel lifecycle + on-chain transactions
  lab02-code              Print the go-perun code to read (interfaces, phases)

Lab 03
  lab03 <scenario>        One scenario: ${LAB03_ALL// /, }
  lab03 all               Run every scenario in turn

Lab 04
  lab04 fast              Cooperative settle (fast path)
  lab04 dispute [secs]    Uncooperative settle, challenge duration defaults to 10
  lab04 all               fast, dispute 10, dispute 30, then compare the timings

Lab 05 (CKB)
  lab05-tools             Install the CKB tooling (apt, Rust target, ckb, ckb-cli). Asks first
  lab05-build             Build contract RISC-V
  lab05-devnet start|stop|status
  lab05 [default|sudt|inspect|dispute]
  lab05-code              Print the contract code to read

Logs go to: $LOGS
EOF
}

cmd=${1:-help}; shift || true
case "$cmd" in
  check)       bash "$CODE/tools/check.sh" ;;
  lab02)       eth_run lab02 run ;;
  lab02-code)  bash "$CODE/tools/read_code.sh" lab02 ;;
  lab03)
    sc=${1:-}
    if [ -z "$sc" ]; then (cd "$CODE/eth" && go run ./cmd/lab03 -h); exit 2; fi
    if [ "$sc" = all ]; then
      for s in $LAB03_ALL; do eth_run lab03 "$s" -scenario "$s"; done
    else
      eth_run lab03 "$sc" -scenario "$sc"
    fi ;;
  lab04)
    m=${1:-dispute}
    case "$m" in
      fast)    eth_run lab04 fast -mode fast ;;
      dispute) eth_run lab04 "dispute${2:-10}" -mode dispute -challenge "${2:-10}" ;;
      all)
        runs=()
        eth_run lab04 fast -mode fast;                  runs+=("$LAST_LOG")
        eth_run lab04 dispute10 -mode dispute -challenge 10; runs+=("$LAST_LOG")
        eth_run lab04 dispute30 -mode dispute -challenge 30; runs+=("$LAST_LOG")
        say "Comparing the 3 runs"
        for f in "${runs[@]}"; do
          echo "--- $(basename "$f")"
          grep -hE "settle took|Total settle time|EVENT|^Total:" "$f"
        done ;;
      *) usage; exit 2 ;;
    esac ;;
  lab05-tools)  bash "$CODE/ckb/tools.sh" ;;
  lab05-build)  bash "$CODE/ckb/build_contracts.sh" ;;
  lab05-devnet) bash "$CODE/ckb/devnet.sh" "${1:-status}" ;;
  lab05)        bash "$CODE/ckb/run.sh" "${1:-default}" ;;
  lab05-code)   bash "$CODE/tools/read_code.sh" lab05 ;;
  help|-h|--help) usage ;;
  *) usage; exit 2 ;;
esac
