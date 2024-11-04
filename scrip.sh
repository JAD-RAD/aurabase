#!/bin/bash

# Configuration
OLD_PATH="github.com/jadd/aurabase"
NEW_PATH="github.com/JAD-RAD/aurabase"

# Update go.mod file
sed -i "s|module $OLD_PATH|module $NEW_PATH|" go.mod

# Find and replace in all .go files
find . -type f -name "*.go" -exec sed -i "s|\"$OLD_PATH|\"$NEW_PATH|g" {} +

# Update the import paths in any go.sum entries (if exists)
if [ -f go.sum ]; then
    sed -i "s|$OLD_PATH|$NEW_PATH|g" go.sum
fi

# Clean go module cache
go clean -modcache

# Tidy up the modules
go mod tidy