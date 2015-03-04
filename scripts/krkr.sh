#!/bin/sh

echo '{
  "piiiiing": "'$(curl -s io:4242/ping)'"
}'
