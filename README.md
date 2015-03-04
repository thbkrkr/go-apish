make build

    docker run --rm -v /home/thb/dev/projects/go/src/krkr:/src -t -i centurylink/golang-builder
    Building github.com/thbkrkr/go-apish

./go-apish &  
    
    [1] 18461
    2015/03/04 14:22:19 Magic server started on port :4242

cat scripts/date.sh 
    
    #!/bin/bash
    echo '{
      "date": '$(date +%s)'
    }'

curl -H 'X-apish-auth:42' localhost:4242/date.sh

    {
      "date": 1425475372
    }
