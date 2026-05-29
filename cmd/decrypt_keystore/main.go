package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	keystoreDir := flag.String("keystore-dir", "./data/geth/keystore", "Path to Geth keystore directory")
	addressFlag := flag.String("address", "", "Ethereum address to decrypt (0x...)")
	password := flag.String("password", "", "Keystore password. Use -password '' for an empty password")
	flag.Parse()

	if *addressFlag == "" {
		fmt.Fprintln(os.Stderr, "missing required -address flag")
		os.Exit(1)
	}
	if !common.IsHexAddress(*addressFlag) {
		fmt.Fprintf(os.Stderr, "invalid Ethereum address: %s\n", *addressFlag)
		os.Exit(1)
	}

	target := strings.ToLower(common.HexToAddress(*addressFlag).Hex()[2:])

	entries, err := os.ReadDir(*keystoreDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read keystore dir: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(*keystoreDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		key, err := keystore.DecryptKey(data, *password)
		if err != nil {
			continue
		}
		got := strings.ToLower(key.Address.Hex()[2:])
		if got != target {
			continue
		}

		privHex := hex.EncodeToString(key.PrivateKey.D.FillBytes(make([]byte, 32)))
		fmt.Printf("Address: %s\n", key.Address.Hex())
		fmt.Printf("Private key (hex, no 0x prefix): %s\n", privHex)
		return
	}

	fmt.Fprintf(os.Stderr, "no decryptable keystore file matched address %s in %s\n", *addressFlag, *keystoreDir)
	os.Exit(1)
}
