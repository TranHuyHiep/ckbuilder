#!/usr/bin/env bash
# Lab 01: checks the server environment and test-builds the modules.
PERUN="${PERUN:-/root/perun-examples}"
CODE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

ok()   { printf '  \033[32m✔\033[0m %-28s %s\n' "$1" "$2"; }
bad()  { printf '  \033[31m✘\033[0m %-28s %s\n' "$1" "$2"; }
ver()  { command -v "$1" >/dev/null && ok "$1" "$($2 2>&1 | head -1)" || bad "$1" "not installed${3:+ ($3)}"; }

echo "== System"
ok "OS" "$(. /etc/os-release; echo "$PRETTY_NAME") $(uname -m)"
ok "CPU / RAM" "$(nproc) cores / $(free -h | awk '/Mem:/{print $2}')"
avail=$(df -BG --output=avail / | tail -1 | tr -dc 0-9)
if [ "$avail" -ge 5 ]; then ok "Disk / free" "${avail} GB"; else bad "Disk / free" "${avail} GB, need ≥ 5 GB"; fi

echo "== Labs 01–04 (Ethereum)"
ver go "go version"
ver node "node --version"
ver npm "npm --version"
[ -d "$PERUN/payment-channel" ] && ok "perun-examples" "$PERUN ($(git -C "$PERUN" log --oneline -1 | cut -c1-50))" || bad "perun-examples" "$PERUN not found"
[ -d "$PERUN/payment-channel/node_modules/hardhat" ] \
  && ok "hardhat" "$(cd "$PERUN/payment-channel" && npx hardhat --version 2>/dev/null)" \
  || bad "hardhat" "not installed, run: cd $PERUN/payment-channel && npm install"

echo "== Lab 05 (CKB)"
ver ckb "ckb --version" "./run.sh lab05-tools"
ver ckb-cli "ckb-cli --version" "./run.sh lab05-tools"
ver tmux "tmux -V" "./run.sh lab05-tools"
ver tmuxp "tmuxp --version" "./run.sh lab05-tools"
ver expect "echo yes" "./run.sh lab05-tools"
ver riscv64-unknown-elf-gcc "riscv64-unknown-elf-gcc --version" "./run.sh lab05-tools"
ver clang-18 "clang-18 --version" "./run.sh lab05-tools"
(cd "$PERUN/payment-channel-ckb/devnet/contracts" && rustup target list --installed 2>/dev/null) | grep -q riscv64imac-unknown-none-elf \
  && ok "rust riscv64imac target" "yes" || bad "rust riscv64imac target" "not installed (./run.sh lab05-tools)"

echo "== Test build"
for m in "$PERUN/payment-channel" "$CODE/eth"; do
  if (cd "$m" && go build ./... 2>/tmp/build.err); then ok "go build" "$m"; else bad "go build" "$m: $(head -3 /tmp/build.err)"; fi
done
if [ -f "$CODE/ckb/go.mod" ]; then
  if (cd "$CODE/ckb" && go build ./... 2>/tmp/build.err); then ok "go build" "$CODE/ckb"; else bad "go build" "$CODE/ckb: $(head -3 /tmp/build.err)"; fi
fi

echo "== Library versions (Lab 01 exercises 1, 2)"
grep -E "go-perun |perun-eth-backend |go-ethereum " "$PERUN/payment-channel/go.mod" | sed 's/^/  payment-channel:     /'
grep -E "go-perun |perun-ckb-backend |ckb-sdk-go" "$PERUN/payment-channel-ckb/go.mod" | sed 's/^/  payment-channel-ckb: /'
echo "== payment-channel-* examples (exercise 3)"
for d in "$PERUN"/payment-channel*; do
  printf '  %-26s %s\n' "$(basename "$d")" "$(head -1 "$d/README.md" 2>/dev/null)"
done
