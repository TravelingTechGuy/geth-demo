package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect(rpcURL string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcURL)
}

func WaitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-timeout:
			return nil, fmt.Errorf("timed out waiting for receipt: %s", txHash.Hex())

		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil && receipt != nil {
				return receipt, nil
			}
		}
	}
}
