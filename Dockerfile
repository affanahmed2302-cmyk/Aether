FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /aether-node ./cmd/aether-node

FROM alpine:3.20
COPY --from=builder /aether-node /aether-node
EXPOSE 7001 7002 7003
ENTRYPOINT ["/aether-node"]
