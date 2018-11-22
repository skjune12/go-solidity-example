package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Config struct {
	Client   ClientConfig   `mapstructure:"client"`
	Contract ContractConfig `mapstructure:"contract"`
}

type ClientConfig struct {
	Url        string `mapstructure:"url"`
	PublicKey  string `mapstructure:"pubKey"`
	PrivateKey string `mapstructure:"PrivKey"`
}

type ContractConfig struct {
	Address string `mapstructure:"address"`
}

var (
	deploy = flag.Bool("deploy", false, "Whether to deploy contract")
	read   = flag.Bool("read", false, "Whether to read data from contract")
	write  = flag.Bool("write", false, "Whether to write data to contract")
)

func LoadConfiguration() *Config {
	// Load Configuration
	viper := viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath("$GOPATH/src/github.com/skjune12/go-solidity-example/contract")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Couldn't load config:", err)
		os.Exit(1)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Couldn't load config:", err)
		os.Exit(1)
	}

	return &config
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "[WARN] Please set some arguments")
		os.Exit(1)
	}

	config := LoadConfiguration()
	flag.Parse()

	client, err := ethclient.Dial(config.Client.Url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to `%s`\n", config.Client.Url)

	// deploy contract if deploy flag is set.
	if *deploy {
		privateKey, err := crypto.HexToECDSA(config.Client.PrivateKey[2:])
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth := bind.NewKeyedTransactor(privateKey)
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0) // in wei
		auth.GasLimit = uint64(300000)
		auth.GasPrice = gasPrice

		input := "1.0"
		address, tx, _, err := DeployStore(auth, client, input)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(address.Hex())
		fmt.Println(tx.Hash().Hex())
	}

	if *read {
		address := common.HexToAddress(config.Contract.Address)
		instance, err := NewStore(address, client)
		if err != nil {
			log.Fatal(err)
		}

		version, err := instance.Version(nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(version)
	}

	if *write {
		privateKey, err := crypto.HexToECDSA(config.Client.PrivateKey[2:])
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public ket to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth := bind.NewKeyedTransactor(privateKey)

		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)
		auth.GasLimit = uint64(300000)
		auth.GasPrice = gasPrice

		// contract address
		address := common.HexToAddress(config.Contract.Address)
		instance, err := NewStore(address, client)

		if err != nil {
			log.Fatal(err)
		}

		key := [32]byte{}
		value := [32]byte{}
		copy(key[:], []byte("hello"))
		copy(value[:], []byte("world"))

		tx, err := instance.SetItem(auth, key, value)

		if err != nil {
			log.Fatal("instance.SetItem", err)
		}

		fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

		result, err := instance.Items(nil, key)
		if err != nil {
			log.Fatal("instance.Item", err)
		}

		fmt.Println(string(result[:]))
	}
}
