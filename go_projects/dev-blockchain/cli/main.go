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

	//add block command
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//holds the block data for the parsed addBlockCmd argument
	addBlockData := addBlockCmd.String("data", "", "Block data")

	//loop over the args and check for the commands
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse((os.Args[2:]))
		if err != nil {
			log.Fatal(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//checks if there's no parsed value..if not provides a description on how to use the command
	if addBlockCmd.Parsed() {
		if *addBlockData == ""{
			addBlockCmd.Usage()
			os.Exit(1)
		}
		//data has been captured store to the chain by dereferencing the parsed argument
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}


func (cli *Cli) addBlock(data string) {
	cli.Bc.AddBlock(data)

	fmt.Println("Success !")
}

func (cli *Cli) printChain() {
	bci := cli.Bc.Iterator()

	//iterates through the block in the chain
	for {
		block, err := bci.Next()

		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("Prev: hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data.Data)
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

//prints the usage of the commands
func (cli *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" addblock -data (BLOCK_DATA) - add a block to the blockchain")
	fmt.Println(" printchain - print all the blocks of the blockchain")
}


//validate that all args are passed
func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}