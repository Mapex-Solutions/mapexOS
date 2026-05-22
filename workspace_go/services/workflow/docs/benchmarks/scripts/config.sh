#!/bin/bash
# Workflow service benchmark configuration

# Service
export SERVICE_URL="http://localhost:5010"
export SERVICE_NAME="workflow"

# Test parameters
export REQUESTS=10000
export CONCURRENCY=50
export DURATION="30s"

# NATS
export NATS_URL="nats://localhost:4222"
export NATS_USERNAME="service"
export NATS_PASSWORD="service_secret"

# Output
export RESULTS_DIR="$(dirname "$0")/../results"
export TAG="${TAG:-$(date +%Y%m%d-%H%M%S)}"
