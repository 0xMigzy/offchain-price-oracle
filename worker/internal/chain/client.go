package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

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
	chainID    *big.Int
}

// NewClient creates a new chain client.
func NewClient(
	ctx context.Context,
	rpcURL string,
	contractAddress string,
	privateKeyHex string,
) (*Client, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid contract address")
	}
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
		chainID:    chainID,
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
	price, ts, err := c.contract.GetLastPrice(opts, pair)
	if err != nil {
		return nil, 0, fmt.Errorf("get latest price: %w", err)
	}
	return price, ts.Uint64(), nil
}

// Address returns the address this client is using to sign transactions.
func (c *Client) Address() common.Address {
	return c.address
}

func (c *Client) NetworkID() uint64 { return c.chainID.Uint64() }

// PreparePrice signs without broadcasting, allowing durable job recovery before
// any externally visible side effect.
func (c *Client) PreparePrice(ctx context.Context, pair string, price *big.Int) (*types.Transaction, error) {
	opts := *c.transactor
	opts.Context = ctx
	opts.NoSend = true
	return c.contract.SubmitPrice(&opts, pair, price)
}

func (c *Client) Broadcast(ctx context.Context, tx *types.Transaction) error {
	return c.eth.SendTransaction(ctx, tx)
}

func (c *Client) Receipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	return c.eth.TransactionReceipt(ctx, hash)
}

func (c *Client) TransactionKnown(ctx context.Context, hash common.Hash) (bool, error) {
	_, _, err := c.eth.TransactionByHash(ctx, hash)
	return err == nil, err
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) { return c.eth.ChainID(ctx) }
func (c *Client) Close()                                        { c.eth.Close() }

// RPCContext bounds individual operations while preserving shutdown cancellation.
func RPCContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 15*time.Second)
}
