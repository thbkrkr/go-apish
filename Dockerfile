FROM alpine:3.2

RUN apk --update add bash && \
    rm -rf /var/cache/apk/*

COPY apish /
COPY api /api

ENTRYPOINT ["/apish"]