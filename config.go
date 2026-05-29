package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RPCURL      string
	PrivateKey1 string
	PrivateKey2 string
}

func LoadConfig() (Config, error) {
	file, err := os.Open(".env")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid .env line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		values[key] = val
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		RPCURL:      values["GETH_RPC_URL"],
		PrivateKey1: values["PRIVATE_KEY_1"],
		PrivateKey2: values["PRIVATE_KEY_2"],
	}

	if cfg.RPCURL == "" || cfg.PrivateKey1 == "" || cfg.PrivateKey2 == "" {
		return Config{}, fmt.Errorf("missing required config values")
	}

	return cfg, nil
}
