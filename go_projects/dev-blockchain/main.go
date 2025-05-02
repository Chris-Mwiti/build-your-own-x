package main

import (
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/blockchain"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/cli"
)

func main(){
	bc := blockchain.NewBlockChain("")
	cli := cli.Cli{
		Bc: bc,
	}
	cli.Run()
}