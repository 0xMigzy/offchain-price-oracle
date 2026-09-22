# Off-Chain Price Oracle

Go services fetch prices from external APIs and submit them on-chain to a Solidity contract, using Kafka for queuing and Redis for caching.

**Status:** Smart contract complete (55 tests, 100% coverage). 
Go chain client submits prices and reads them back, with transaction simulation and retry logic, verified end-to-end against a local Anvil node.
Go services, Docker and CI in progress.
