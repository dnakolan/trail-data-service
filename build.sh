#!/bin/bash

# Build dependencies image
echo "Building dependencies image..."
docker build -t trail-data-service-deps -f Dockerfile.deps .

# Build main application image
echo "Building application image..."
docker build -t trail-data-service -f Dockerfile .

echo "Build complete!" 