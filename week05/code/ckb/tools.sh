#!/usr/bin/env bash
# Lab 05: installs the CKB tooling on the Ubuntu 26.04 arm64 server. Safe to re-run.
set -euo pipefail

CKB_VER=v0.210.0
CKBCLI_VER=v1.16.0
ARCH=$(uname -m)   # aarch64 on this server
case "$ARCH" in
  aarch64|x86_64) ;;
  *) echo "No prebuilt ckb release for architecture $ARCH"; exit 1 ;;
esac

say() { printf '\n\033[1;36m>>> %s\033[0m\n' "$*"; }

cat <<EOF
This script will:
  1. apt install: tmux tmuxp expect python-is-python3 gcc-riscv64-unknown-elf
     binutils-riscv64-unknown-elf libc6-dev-riscv64-cross linux-libc-dev-riscv64-cross clang-18 llvm-18
  2. rustup toolchain 1.85.0 + target riscv64imac-unknown-none-elf, pinned only for the contracts directory
     (newer Rust cannot build the contracts; stable and 1.95.0 used by /root/fiber are untouched)
  3. Install ckb $CKB_VER to /usr/local/bin/ckb
  4. Replace ckb-cli $(ckb-cli --version 2>/dev/null | awk '{print $2}') with $CKBCLI_VER
     (backing up /usr/local/bin/ckb-cli and ~/.ckb-cli first)
EOF
read -r -p "Continue? [y/N] " ans
[[ "$ans" =~ ^[yY]$ ]] || { echo "Cancelled."; exit 0; }

avail=$(df -BG --output=avail / | tail -1 | tr -dc 0-9)
[ "$avail" -lt 5 ] && { echo "Only ${avail} GB free on /, need ≥ 5 GB"; exit 1; }

say "1/4 apt install"
apt-get update
apt-get install -y tmux tmuxp expect python-is-python3 \
  gcc-riscv64-unknown-elf binutils-riscv64-unknown-elf \
  libc6-dev-riscv64-cross linux-libc-dev-riscv64-cross \
  clang-18 llvm-18

# The contracts' C code is compiled by clang-18 for soft-float riscv64 (ABI lp64), but Ubuntu's
# libc6-dev-riscv64-cross only ships headers for lp64d. stubs-lp64.h only lists unimplemented
# functions, so the lp64d copy is a safe stand-in. Adds two files, changes nothing existing.
G=/usr/riscv64-linux-gnu/include/gnu
for f in stubs lib-names; do
  [ -e "$G/$f-lp64.h" ] || cp "$G/$f-lp64d.h" "$G/$f-lp64.h"
done

say "2/4 Rust 1.85.0 + RISC-V target for the contracts"
# The contracts build with -C target-feature=-a (CKB-VM has no atomics). Rust ≥ 1.9x no longer
# compiles the atomics in the `bytes` crate under that flag ("Cannot select: AtomicLoadAdd"),
# so pin the version the perun-ckb-contract README recommends, for that directory only.
CONTRACTS="${PERUN:-/root/perun-examples}/payment-channel-ckb/devnet/contracts"
rustup toolchain install 1.85.0 --profile minimal
rustup target add riscv64imac-unknown-none-elf --toolchain 1.85.0
(cd "$CONTRACTS" && rustup override set 1.85.0)

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

say "3/4 ckb $CKB_VER"
if ckb --version 2>/dev/null | grep -q "${CKB_VER#v}"; then
  echo "Already installed: $(ckb --version)"
else
  f="ckb_${CKB_VER}_${ARCH}-unknown-linux-gnu"
  wget -q --show-progress -O "$TMP/ckb.tgz" "https://github.com/nervosnetwork/ckb/releases/download/$CKB_VER/$f.tar.gz"
  tar xzf "$TMP/ckb.tgz" -C "$TMP"
  install -m 755 "$TMP/$f/ckb" /usr/local/bin/ckb   # take only ckb, skip the bundled ckb-cli 2.0
fi

say "4/4 ckb-cli $CKBCLI_VER"
if ckb-cli --version 2>/dev/null | grep -q "${CKBCLI_VER#v}"; then
  echo "Already installed: $(ckb-cli --version)"
else
  if [ -x /usr/local/bin/ckb-cli ]; then
    old=$(ckb-cli --version | awk '{print $2}')
    cp -n /usr/local/bin/ckb-cli "/usr/local/bin/ckb-cli-$old.bak"
    echo "Backed up old binary: /usr/local/bin/ckb-cli-$old.bak"
  fi
  if [ -d ~/.ckb-cli ] && [ ! -d ~/.ckb-cli.bak ]; then
    cp -a ~/.ckb-cli ~/.ckb-cli.bak
    echo "Backed up keystore: ~/.ckb-cli.bak"
  fi
  wget -q --show-progress -O "$TMP/cli.tgz" "https://github.com/nervosnetwork/ckb-cli/releases/download/$CKBCLI_VER/ckb-cli_${CKBCLI_VER}_${ARCH}-unknown-linux-gnu.tar.gz"
  mkdir "$TMP/cli" && tar xzf "$TMP/cli.tgz" -C "$TMP/cli"
  install -m 755 "$(find "$TMP/cli" -name ckb-cli -type f | head -1)" /usr/local/bin/ckb-cli
fi

say "Result"
ckb --version
ckb-cli --version
tmux -V; tmuxp --version; python --version
riscv64-unknown-elf-gcc --version | head -1
clang-18 --version | head -1
(cd "$CONTRACTS" && rustc --version && rustup target list --installed | grep riscv)
echo
echo "Done. Next step: ./run.sh lab05-build"
