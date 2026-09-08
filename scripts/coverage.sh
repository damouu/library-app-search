#!/bin/bash

set -e

go test \
  ./internal/config \
  ./internal/elasticsearch \
  ./internal/health \
  -coverprofile=coverage.out \
  -covermode=atomic

go tool cover -func=coverage.out