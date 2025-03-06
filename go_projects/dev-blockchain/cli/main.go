package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/blockchain"
)


type Cli struct {
	Bc *blockchain.Blockchain
}

func (cli *Cli) Run(){
	//validate the cli arguments
	cli.validateArgs()

	//blockchain commands
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createChainCmd := flag.NewFlagSet("createchain", flag.ExitOnError)


	//used to hold the address of the newly created chain	
	chainAddress := createChainCmd.String("address", "", "Chain address")

	//loop over the args and check for the commands
	switch os.Args[1] {
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

	case "createchain":
		err := createChainCmd.Parse((os.Args[2:]))
		if err != nil {
			log.Fatal(err)
		}
	
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createChainCmd.Parsed() {
		if *chainAddress == ""{
			createChainCmd.Usage()
			os.Exit(1)
		}
		cli.createChain(*chainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}



//proxy func to print chains in the blockchain
func (cli *Cli) printChain() {
	bci := cli.Bc.Iterator()

	//iterates through the block in the chain
	for {
		block, err := bci.Next()

		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("Prev: hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		//checks if we have reached the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

//creates a new chain
func (cli *Cli) createChain(address string){ 
	chain := blockchain.BlockChainWithDb(address)
	defer chain.Db.Close()
	fmt.Println("Completed creating the chain!")
}

func (cli *Cli) getBalance(address string){
	bc := blockchain.BlockChainWithDb(address);
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUnspentTxo(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s : %d\n", address, balance)
}

//prints the usage of the commands
func (cli *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" printchain - print all the blocks of the blockchain")
	fmt.Println("createchain - creates a chain if none exists")
}


//validate that all args are passed
func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}