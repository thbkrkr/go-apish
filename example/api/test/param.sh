#!/bin/sh
set -eu

# Build JSON with jq so the parameter is always safely quoted/escaped,
# rather than interpolating it straight into the output.
jq -n --arg param "${1:-}" '{param: $param}'
