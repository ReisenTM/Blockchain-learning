package part_v

import "testing"

func TestWallet_GetAddress(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	t.Log(string(address))
}
