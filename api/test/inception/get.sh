#!/bin/sh

c1="curl -s io:4242"
c2="curl -s io:4242/api/time/date"
c3='curl -s -u zuperadmin:42 io:4242/api/time/date'
c4='curl -s -u zuperadmin:42 io:4242/api/test/ping'

w='{"status":"%{http_code}","time":"%{time_total}"}'

slurp () { jq -s .; }

echo '{
  "response1": '$($c1 -w $w | slurp)',
  "response2": '$($c2 -w $w | slurp)',
  "response3": '$($c3 -w $w | slurp)',
  "response4": '$($c4 -w $w | slurp)'
}'
