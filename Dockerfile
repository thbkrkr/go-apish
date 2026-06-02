# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG GIT_COMMIT=undefined
ARG BUILD_DATE=undefined
RUN CGO_ENABLED=0 go build \
    -ldflags "-X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${BUILD_DATE}" \
    -o /go-apish ./app

# Runtime stage
FROM alpine:3.20

RUN apk --no-cache add bash jq curl
COPY --from=build /go-apish /go-apish

EXPOSE 4242
ENTRYPOINT ["/go-apish"]
