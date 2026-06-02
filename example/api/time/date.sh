#!/bin/sh
set -eu

# Build JSON with jq so values are always correctly typed and escaped.
jq -n \
  --argjson date "$(date +%s)" \
  --arg human_date "$(date)" \
  '{date: $date, human_date: $human_date}'
