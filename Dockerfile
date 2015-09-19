FROM alpine:3.2

RUN apk --update add bash jq curl && \
    rm -rf /var/cache/apk/*

COPY go-apish /

CMD ["/go-apish"]
