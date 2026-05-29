package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	client, err := Connect(cfg.RPCURL)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("chain ID failed: %v", err)
	}

	walletA, err := NewWalletFromHex(cfg.PrivateKey1)
	if err != nil {
		log.Fatalf("wallet A failed: %v", err)
	}

	walletB, err := NewWalletFromHex(cfg.PrivateKey2)
	if err != nil {
		log.Fatalf("wallet B failed: %v", err)
	}

	fmt.Println("Wallet A:", walletA.Address().Hex())
	fmt.Println("Wallet B:", walletB.Address().Hex())

	balA1, err := walletA.Balance(context.Background(), client)
	if err != nil {
		log.Fatalf("wallet A balance failed: %v", err)
	}

	balB1, err := walletB.Balance(context.Background(), client)
	if err != nil {
		log.Fatalf("wallet B balance failed: %v", err)
	}

	fmt.Println("Before transfer")
	fmt.Println("Wallet A:", WeiToEthString(balA1), "ETH")
	fmt.Println("Wallet B:", WeiToEthString(balB1), "ETH")

	tx, err := TransferETH(client, chainID, walletA, walletB.Address(), EthToWei("1"))
	if err != nil {
		log.Fatalf("transfer failed: %v", err)
	}

	fmt.Println("Submitted tx:", tx.Hash().Hex())

	receipt, err := WaitForReceipt(context.Background(), client, tx.Hash())
	if err != nil {
		log.Fatalf("failed waiting for receipt: %v", err)
	}

	fmt.Println("Mined in block:", receipt.BlockNumber.Uint64())

	balA2, err := walletA.Balance(context.Background(), client)
	if err != nil {
		log.Fatalf("wallet A balance failed: %v", err)
	}

	balB2, err := walletB.Balance(context.Background(), client)
	if err != nil {
		log.Fatalf("wallet B balance failed: %v", err)
	}

	fmt.Println("After transfer")
	fmt.Println("Wallet A:", WeiToEthString(balA2), "ETH")
	fmt.Println("Wallet B:", WeiToEthString(balB2), "ETH")
}
