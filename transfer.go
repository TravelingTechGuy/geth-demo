package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransferETH(
	client *ethclient.Client,
	chainID *big.Int,
	from *Wallet,
	to common.Address,
	valueWei *big.Int,
) (*types.Transaction, error) {
	ctx := context.Background()

	nonce, err := client.PendingNonceAt(ctx, from.Address())
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}

	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, to, valueWei, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("send tx: %w", err)
	}

	return signedTx, nil
}
