#!/usr/bin/env bash
# Lab 05: builds the Perun CKB contracts (RISC-V). make dev rebuilds them too, but
# building separately first makes errors easier to see.
set -euo pipefail
PERUN="${PERUN:-/root/perun-examples}"
C="$PERUN/payment-channel-ckb/devnet/contracts"

for t in clang-18 riscv64-unknown-elf-gcc; do
  command -v $t >/dev/null || { echo "Missing $t. Run: ./run.sh lab05-tools"; exit 1; }
done
rustup target list --installed | grep -q riscv64imac-unknown-none-elf \
  || { echo "Missing Rust target riscv64imac. Run: ./run.sh lab05-tools"; exit 1; }

cd "$C"
start=$(date +%s)
# shellcheck disable=SC1091
source ./setup_env.sh build
make build
echo
echo "Build finished after $(( $(date +%s) - start )) seconds:"
ls -la build/release
du -sh target
