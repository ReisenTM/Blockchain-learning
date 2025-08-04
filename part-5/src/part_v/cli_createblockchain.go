package part_v

import "fmt"

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()
	fmt.Println("Done!")
}
