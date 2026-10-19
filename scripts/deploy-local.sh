#!/bin/bash
set -euo pipefail

: "${PRIVATE_KEY:?set PRIVATE_KEY to the deployer private key}"

echo "==> bringing up infra"
docker compose up -d anvil redis kafka

echo "==> waiting for anvil"
until curl -s http://localhost:8545 > /dev/null; do sleep 1; done

echo "==> deploying contract"
cd contracts
DEPLOYED=$(forge script script/Deploy.s.sol \
    --rpc-url http://localhost:8545 \
    --private-key "$PRIVATE_KEY" \
    --broadcast 2>&1 | grep "deployed at" | awk '{print $NF}')
cd ..

echo "==> deployed at $DEPLOYED"
printf 'CONTRACT_ADDRESS=%s\nPRIVATE_KEY=%s\n' "$DEPLOYED" "$PRIVATE_KEY" > .env
echo "==> running worker"
docker compose up worker
