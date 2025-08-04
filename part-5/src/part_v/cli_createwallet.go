package part_v

import (
	"fmt"
)

func (cli *CLI) createWallet() {
	wallet := NewWallet()
	addr := wallet.GetAddress()
	fmt.Println("New Wallet:", string(addr))
}
