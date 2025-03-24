#!/bin/bash
set -e  # Exit on error

docker compose version

docker compose --profile test build
docker compose --profile test run test