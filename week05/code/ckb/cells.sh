#!/usr/bin/env bash
# Lab 05: lists the live Perun cells on the devnet (channel cell, funds cells).
set -uo pipefail
PERUN="${PERUN:-/root/perun-examples}"
MIG="$PERUN/payment-channel-ckb/devnet/contracts/migrations/dev"
RPC=http://127.0.0.1:8114

hash_of() { jq -r --arg n "$1" '.cell_recipes[] | select(.name == $n) | .data_hash' "$MIG"/*.json | head -1; }

PCTS=$(hash_of pcts); PCLS=$(hash_of pcls); PFLS=$(hash_of pfls); SUDT=$(hash_of sudt)
echo "code_hash (hash_type = data1, from $MIG):"
printf '  %-5s %s\n' PCTS "$PCTS" PCLS "$PCLS" PFLS "$PFLS" SUDT "$SUDT"

name_of() {
  case "$1" in
    "$PCTS") echo PCTS ;; "$PCLS") echo PCLS ;; "$PFLS") echo PFLS ;; "$SUDT") echo SUDT ;;
    null|"") echo "-" ;; *) echo "${1:0:10}…" ;;
  esac
}

# query <title> <code_hash> <lock|type>
query() {
  local body
  body=$(jq -n --arg h "$2" --arg t "$3" '{id:1,jsonrpc:"2.0",method:"get_cells",
    params:[{script:{code_hash:$h,hash_type:"data1",args:"0x"},script_type:$t,script_search_mode:"prefix"},"asc","0x20"]}')
  echo
  echo "== $1"
  curl -s -X POST -H 'content-type: application/json' -d "$body" $RPC \
    | jq -r '.result.objects[] | [.out_point.tx_hash, .out_point.index, .output.capacity,
             .output.lock.code_hash, (.output.type.code_hash // "null"), ((.output_data | length) / 2 - 1)] | @tsv' \
    | while IFS=$'\t' read -r tx idx cap lock type dlen; do
        printf '  out_point %s…:%d  capacity %s CKB  lock=%s  type=%s  data=%d byte\n' \
          "${tx:0:14}" "$((idx))" "$(python3 -c "print($((cap))/1e8)")" "$(name_of "$lock")" "$(name_of "$type")" "$dlen"
      done
}

query "Channel cell (type script = PCTS)" "$PCTS" type
query "Funds cell (lock script = PFLS)" "$PFLS" lock
