#!/bin/sh

# Kill process if already running on port
lsof -i:3000 -Fp | head -n 1 | sed 's/^p//' | xargs kill

# Start using nodemon
nodemon --exec go run main.go --signal SIGTERM

