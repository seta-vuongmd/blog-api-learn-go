#!/bin/sh
set -e
echo "Current directory: $(pwd)"
echo "Listing files in current directory:"
ls -la
echo "Running application..."
exec ./blog-api-learn-go