FROM alpine:3.7

RUN apk --no-cache add bash jq curl
COPY go-apish /go-apish
CMD ["/go-apish"]
