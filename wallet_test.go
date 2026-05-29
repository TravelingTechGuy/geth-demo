package main

import "testing"

func TestNewWalletFromHex(t *testing.T) {
	w, err := NewWalletFromHex("0xc87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := w.Address().Hex(); got != "0x627306090abaB3A6e1400e9345bC60c78a8BEf57" {
		t.Fatalf("unexpected address: %s", got)
	}
}
