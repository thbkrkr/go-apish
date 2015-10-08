#!/bin/sh
set -eu

c1="curl -s -m 1 172.17.42.1"
c2="curl -s -m 1 172.17.42.1/api/time/date"
c3='curl -s -m 1 -u zuperadmin:42 172.17.42.1/api/time/date'
c4='curl -s -m 1 -u zuperadmin:42 172.17.42.1/api/test/param?q=hello'

w='{"status":"%{http_code}","time":"%{time_total}"}'

slurp () { jq -s .; }

echo '{
  "1": '$($c1 -w $w | slurp)',
  "2": '$($c2 -w $w | slurp)',
  "3": '$($c3 -w $w | slurp)',
  "4": '$($c4 -w $w | slurp)'
}'
