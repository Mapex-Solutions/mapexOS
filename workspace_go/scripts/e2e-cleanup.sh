#!/bin/bash

# E2E Test Cleanup Script
# Limpa MongoDB, Redis e NATS para garantir ambiente limpo antes dos testes E2E

set -e

echo "🧹 Starting E2E environment cleanup..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Clean MongoDB
echo -e "${YELLOW}📦 Cleaning MongoDB...${NC}"
mongosh --quiet --eval "
  const databases = ['mapexos', 'assets', 'router', 'http_gateway'];
  databases.forEach(dbName => {
    const dbInstance = db.getSiblingDB(dbName);
    const collections = dbInstance.getCollectionNames();
    collections.forEach(collection => {
      dbInstance[collection].deleteMany({});
      print(\`  ✓ Cleaned collection: \${dbName}.\${collection}\`);
    });
  });
" || echo -e "${RED}  ✗ MongoDB cleanup failed${NC}"

echo ""

# 2. Clean Redis
echo -e "${YELLOW}🔴 Cleaning Redis...${NC}"
redis-cli FLUSHALL > /dev/null && echo -e "${GREEN}  ✓ Redis FLUSHALL executed${NC}" || echo -e "${RED}  ✗ Redis cleanup failed${NC}"

echo ""

# 3. Clean NATS Streams
echo -e "${YELLOW}📨 Cleaning NATS streams...${NC}"

# List of streams to purge
STREAMS=("ROUTE-GROUPS" "ASSET-UPDATED-ROUTING" "ASSET-UPDATED-SCRIPTS")

for stream in "${STREAMS[@]}"; do
  if nats stream info "$stream" &> /dev/null; then
    nats stream purge "$stream" --force > /dev/null 2>&1 && \
      echo -e "${GREEN}  ✓ Purged stream: $stream${NC}" || \
      echo -e "${YELLOW}  ⚠ Stream $stream not found or already empty${NC}"
  else
    echo -e "${YELLOW}  ⚠ Stream $stream does not exist${NC}"
  fi
done

echo ""
echo -e "${GREEN}✅ E2E environment cleanup completed!${NC}"
echo ""
echo "Environment is ready for E2E tests:"
echo "  • MongoDB: All collections cleared"
echo "  • Redis: All keys flushed"
echo "  • NATS: All streams purged"
echo ""
