package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func NewWalletFromHex(hexKey string) (*Wallet, error) {
	hexKey = strings.TrimPrefix(strings.TrimSpace(hexKey), "0x")

	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Wallet{
		privateKey: privateKey,
		address:    address,
	}, nil
}

func (w *Wallet) Address() common.Address {
	return w.address
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) Balance(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	return client.BalanceAt(ctx, w.address, nil)
}

func EthToWei(eth string) *big.Int {
	f, ok := new(big.Float).SetString(eth)
	if !ok {
		panic(fmt.Sprintf("invalid ETH amount: %s", eth))
	}

	weiFloat := new(big.Float).Mul(f, big.NewFloat(1e18))
	wei := new(big.Int)
	weiFloat.Int(wei)
	return wei
}

func WeiToEthString(wei *big.Int) string {
	f := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return f.Text('f', 6)
}
