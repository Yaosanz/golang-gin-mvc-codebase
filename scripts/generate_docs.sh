#!/bin/bash

# exec: bash scripts/generate_docs.sh
echo "Generating Swagger documentation..."
swag init -g cmd/app/main.go
echo "Swagger documentation generated!"