// test-submit is a manual smoke-test tool: it submits one price to a
// deployed OracleConsumer contract and reads it back to confirm the round
// trip works end-to-end against a live RPC endpoint (e.g. a local anvil node).
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/0xmigzy/capstone-worker/bindings/oracleconsumer"
	"github.com/0xmigzy/capstone-worker/internal/chain"
)

// anvilDefaultKey is anvil's well-known account #0 private key. It is public
// knowledge shipped with foundry and only ever holds test funds.
const anvilDefaultKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func main() {
	rpcURL := flag.String("rpc", "http://127.0.0.1:8545", "EVM JSON-RPC endpoint")
	contractAddr := flag.String("contract", "", "OracleConsumer contract address (required)")
	privateKey := flag.String("key", anvilDefaultKey, "hex private key of an authorized submitter")
	pair := flag.String("pair", "ETH/USD", "price pair identifier")
	price := flag.Int64("price", 300000000000, "price to submit (integer, e.g. 8-decimal fixed point)")
	timeout := flag.Duration("timeout", 30*time.Second, "overall timeout")
	flag.Parse()

	if *contractAddr == "" {
		log.Fatal("missing required -contract flag")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client, err := chain.Dial(ctx, *rpcURL)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer client.Close()

	auth, err := client.NewTransactor(*privateKey)
	if err != nil {
		log.Fatalf("build transactor: %v", err)
	}

	oracle, err := oracleconsumer.NewOracleconsumer(common.HexToAddress(*contractAddr), client)
	if err != nil {
		log.Fatalf("bind contract: %v", err)
	}

	fmt.Printf("submitting %s = %d via %s...\n", *pair, *price, auth.From)

	tx, err := oracle.SubmitPrice(auth, *pair, big.NewInt(*price))
	if err != nil {
		log.Fatalf("submit price: %v", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("wait mined: %v", err)
	}
	if receipt.Status != 1 {
		log.Fatalf("transaction reverted (tx %s)", tx.Hash())
	}

	fmt.Printf("mined in block %d (tx %s, gas used %d)\n", receipt.BlockNumber, tx.Hash(), receipt.GasUsed)

	gotPrice, gotTimestamp, err := oracle.GetLastPrice(&bind.CallOpts{Context: ctx}, *pair)
	if err != nil {
		log.Fatalf("get last price: %v", err)
	}

	fmt.Printf("readback: %s -> price=%s timestamp=%s\n", *pair, gotPrice, gotTimestamp)
}
