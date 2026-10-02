#!/usr/bin/env bash
# Simple helper to start a 3-node local cluster for development.

set -e

echo "Starting Aether 3-node cluster..."

go run ./cmd/aether-node --id node1 --port 7001 &
go run ./cmd/aether-node --id node2 --port 7002 &
go run ./cmd/aether-node --id node3 --port 7003 &

echo "Nodes started in background. Use 'pkill -f aether-node' to stop."
