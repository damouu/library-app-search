#!/bin/bash

set -e

go test \
  ./internal/config \
  ./internal/opensearch \
  ./internal/health \
  -coverprofile=coverage.out \
  -covermode=atomic

go tool cover -func=coverage.out