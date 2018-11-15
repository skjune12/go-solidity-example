package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/k0kubun/pp"
)

func main() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("We have a connection")
	_ = client

	address1 := common.HexToAddress("0xd2844e024c5bb23ed8118baeba4a0d7a0e4877d4")

	fmt.Println("address.Hex():", address1.Hex())
	fmt.Println("address.Hash():", address1.Hash().Hex())
	fmt.Println("address.Bytes():", address1.Bytes())

	// account balances
	account := common.HexToAddress("0xffe98770fb5686e71064d9d7ca51c9f86ee200be")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)

	blockNumber := big.NewInt(0)
	balance, err = client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	// generating new wallets
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	pp.Println(hexutil.Encode(privateKeyBytes)[2:]) // omit "0x"

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pp.Println(hexutil.Encode(publicKeyBytes)[4:]) // omit "0x04"

	address2 := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address:", address2)

	fmt.Println(balance)

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
