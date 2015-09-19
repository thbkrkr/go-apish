#!/bin/sh
set -eu

echo '{
  "date": '$(date +%s)',
  "human_date": "'$(date)'"
}'

