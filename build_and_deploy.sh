#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_ROOT=$(pwd)
PROJECT_NAME=$(basename "$PROJECT_ROOT")
BINARY_NAME="$PROJECT_NAME"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Building $PROJECT_NAME${NC}"
echo -e "${GREEN}========================================${NC}"

# Validate main.go exists
if [ ! -f "main.go" ]; then
    echo -e "${RED}Error: main.go not found in project root${NC}"
    exit 1
fi

echo -e "${YELLOW}→ Tidying Go modules${NC}"
go mod tidy

echo -e "${YELLOW}→ Building binary${NC}"
go build -o "$BINARY_NAME" main.go

echo -e "${GREEN}✓ Build successful${NC}"

# PM2 logic
if ! command -v pm2 &> /dev/null; then
    echo -e "${RED}Error: PM2 not installed${NC}"
    echo -e "${YELLOW}Install with: npm install -g pm2${NC}"
    exit 1
fi

echo -e "${YELLOW}→ Starting / Restarting PM2 process${NC}"

if pm2 describe "$PROJECT_NAME" &> /dev/null; then
    pm2 restart "$PROJECT_NAME" --update-env
else
    pm2 start "$PROJECT_ROOT/$BINARY_NAME" --name "$PROJECT_NAME" --cwd "$PROJECT_ROOT"
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Done.${NC}"
echo -e "${GREEN}Binary: $PROJECT_ROOT/$BINARY_NAME${NC}"
echo -e "${GREEN}========================================${NC}"