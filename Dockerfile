FROM golang:1.13.7-alpine3.11 AS builder

WORKDIR /go/build
COPY . .
RUN go build -o iv_tracker cmd/main.go

FROM alpine:3.11.3

WORKDIR /go/run
COPY --from=builder /go/build/iv_tracker iv_tracker
ENTRYPOINT ["./iv_tracker"]