make build run test stop 
    
    docker run --rm -v /home/thb/dev/projects/go/src/krkr:/src -t -i centurylink/golang-builder
    Building github.com/thbkrkr/go-apish
    
    ./go-apish -port=1234 -apiKeyHeader=X-pof-auth -apiKey=E7D9EJJD87EH7ED87H &
    2015/03/04 15:01:02 Magic server started on port :1234

    curl -H 'X-pof-auth:E7D9EJJD87EH7ED87H' localhost:1234/date.sh
    {
      "date": 1425477662
    }

    pkill go-apish
