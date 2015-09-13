#!/bin/bash
set -eu

here=$(dirname $(readlink -f $0))

dateFile=$here/../_static/date.html

date >> $dateFile
echo '<br>' >> $dateFile

echo '{
  "date": '$(date +%s)',
  "human_date": "'$(date)'"
}'

