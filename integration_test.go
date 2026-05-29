package main

import (
	"context"
	"testing"
)

func TestIntegrationTransferETH(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	client, err := Connect(cfg.RPCURL)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatalf("chain ID failed: %v", err)
	}

	ctx := context.Background()

	walletA, err := NewWalletFromHex(cfg.PrivateKey1)
	if err != nil {
		t.Fatalf("wallet A failed: %v", err)
	}

	walletB, err := NewWalletFromHex(cfg.PrivateKey2)
	if err != nil {
		t.Fatalf("wallet B failed: %v", err)
	}

	beforeA, err := walletA.Balance(ctx, client)
	if err != nil {
		t.Fatalf("balance A before failed: %v", err)
	}

	beforeB, err := walletB.Balance(ctx, client)
	if err != nil {
		t.Fatalf("balance B before failed: %v", err)
	}

	tx, err := TransferETH(client, chainID, walletA, walletB.Address(), EthToWei("0.5"))
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}

	if _, err := WaitForReceipt(ctx, client, tx.Hash()); err != nil {
		t.Fatalf("receipt wait failed: %v", err)
	}

	afterA, err := walletA.Balance(ctx, client)
	if err != nil {
		t.Fatalf("balance A after failed: %v", err)
	}

	afterB, err := walletB.Balance(ctx, client)
	if err != nil {
		t.Fatalf("balance B after failed: %v", err)
	}

	if afterA.Cmp(beforeA) >= 0 {
		t.Fatalf("expected wallet A balance to decrease")
	}

	if afterB.Cmp(beforeB) <= 0 {
		t.Fatalf("expected wallet B balance to increase")
	}
}
