#!/bin/bash -eu

IN="$(cat /dev/stdin)"

echo '{
  "jackpot": '$(jq .o <<< $IN)'
}'