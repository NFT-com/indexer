package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/d7b15235a515483490a5b89644221a71")
	if err != nil {
		log.Fatalln(err)
	}

	a, _ := client.CodeAt(context.Background(), common.HexToAddress("0x06012c8cf97bead5deae237070f9587f8e7a266d"), nil)

	fmt.Println(hex.EncodeToString(a))
}
