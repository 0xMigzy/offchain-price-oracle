package testsubmit

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/chain"
	"github.com/0xmigzy/capstone-worker/internal/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcURL := os.Getenv("RPC_URL")
	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")

	chainClient, err := chain.NewClient(ctx, rpcURL, contractAddr, privateKey)
	if err != nil {
		log.Fatalf("create chain client: %v", err)
	}

	oracle := client.NewOracleClient(chainClient)

	// Submit a price
	txHash, err := oracle.Submit(ctx, client.PriceUpdate{
		Pair:  "ETH/USD",
		Price: big.NewInt(3500e8),
	})
	if err != nil {
		log.Fatalf("submit price: %v", err)
	}
	log.Printf("submitted: tx=%s", txHash)

	// Wait briefly for inclusion
	time.Sleep(3 * time.Second)

	// Read it back
	price, ts, err := oracle.GetLatest(ctx, "ETH/USD")
	if err != nil {
		log.Fatalf("get latest: %v", err)
	}
	log.Printf("latest: pair=ETH/USD price=%s timestamp=%d", price.String(), ts)
}
