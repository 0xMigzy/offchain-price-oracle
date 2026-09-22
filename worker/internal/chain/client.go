package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/0xmigzy/capstone-worker/bindings/oracleconsumer"
)

// Client wraps the Ethereum client and contract bindings.
type Client struct {
	eth        *ethclient.Client
	contract   *oracleconsumer.Oracleconsumer
	transactor *bind.TransactOpts
	address    common.Address
}

// NewClient creates a new chain client.
func NewClient(
	ctx context.Context,
	rpcURL string,
	contractAddress string,
	privateKeyHex string,
) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial RPC: %w", err)
	}

	chainID, err := eth.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key conversion failed")
	}
	address := crypto.PubkeyToAddress(*publicKey)

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	contract, err := oracleconsumer.NewOracleconsumer(common.HexToAddress(contractAddress), eth)
	if err != nil {
		return nil, fmt.Errorf("bind contract: %w", err)
	}

	return &Client{
		eth:        eth,
		contract:   contract,
		transactor: transactor,
		address:    address,
	}, nil
}

// SubmitPrice submits a price to the contract. Returns the transaction hash.
func (c *Client) SubmitPrice(ctx context.Context, pair string, price *big.Int) (string, error) {
	opts := *c.transactor // copy
	opts.Context = ctx

	tx, err := c.contract.SubmitPrice(&opts, pair, price)
	if err != nil {
		return "", fmt.Errorf("submit price: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetLatestPrice reads the latest price for a pair.
func (c *Client) GetLatestPrice(ctx context.Context, pair string) (*big.Int, uint64, error) {
	opts := &bind.CallOpts{Context: ctx}
	price, ts, err := c.contract.GetLatestPrice(opts, pair)
	if err != nil {
		return nil, 0, fmt.Errorf("get latest price: %w", err)
	}
	return price, ts.Uint64(), nil
}

// Address returns the address this client is using to sign transactions.
func (c *Client) Address() common.Address {
	return c.address
}
