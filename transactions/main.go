package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("We have a connection")

	// querying blocks
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("blockNumber:", header.Number.String())

	block, err := client.BlockByNumber(context.Background(), big.NewInt(5671746))

	fmt.Println("blockNumber:", block.Number().Uint64())
	fmt.Println("blockTime:", block.Time().Uint64())
	fmt.Println("blockHash:", block.Hash().Hex())
	fmt.Println(len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count)
}
