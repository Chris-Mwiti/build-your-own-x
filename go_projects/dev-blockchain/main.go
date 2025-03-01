package main

import (
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/cli"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/blockchain"
)



func main(){
	chain := blockchain.NewBlockchain()

	//closes the connection
	defer chain.Db.Close()
	
	chainCli := cli.Cli{
		Bc: chain,
	}
	chainCli.Run()

}