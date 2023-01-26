FROM golang:1.19.5-alpine3.17 AS builder

WORKDIR /go/build
COPY . .
RUN go build -o iv_tracker cmd/main.go

FROM alpine:3.17.1

WORKDIR /go/run
COPY db/migrations db/migrations
COPY --from=builder /go/build/iv_tracker iv_tracker
ENTRYPOINT ["./iv_tracker"]