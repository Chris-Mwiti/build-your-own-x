package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/blockchain"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/wallets"
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
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendBalanceCmd := flag.NewFlagSet("send", flag.ExitOnError)


	//used to hold the address of the newly created chain	
	chainAddress := createChainCmd.String("address", "", "Chain address")

	//stores the get balance address
	balanceAddress := getBalanceCmd.String("address", "", "wallet address")

	//stores the from, to and amount to be sent over the network
	senderAddress := sendBalanceCmd.String("from", "", " sender wallet address")
	receiverAddress := sendBalanceCmd.String("to", "", " receiver wallet address")
	amountToSend := sendBalanceCmd.Int("amount", 0, "amount to be sent")


	//wallets commands
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	//

	//loop over the args and check for the commands and their subsets are already parsed
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
	
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

	case "send":
		err := sendBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}		
	case "createwallet":
		err := createChainCmd.Parse(os.Args[2:])
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

	if getBalanceCmd.Parsed() {
		if *balanceAddress == ""{
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*balanceAddress)
	}	
	if sendBalanceCmd.Parsed() {
		if *senderAddress == "" || *receiverAddress == "" || *amountToSend <= 0 {
			sendBalanceCmd.Usage()
			os.Exit(1)
	
		}
		cli.send(*senderAddress, *receiverAddress, *amountToSend)
	}

	if createWalletCmd.Parsed(){
		cli.createWallet()
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
	chain := blockchain.CreateBlockchain(address)
	defer chain.Db.Close()
	fmt.Println("Completed creating the chain!")
}

func (cli *Cli) getBalance(address string){
	bc := blockchain.NewBlockChain(address);
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUnspentTxo([]byte(address))

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s : %d\n", address, balance)
}


func (cli *Cli) send(from, to string, amount int){
	//initialize the blockchain
	bc := blockchain.NewBlockChain(from)

	defer bc.Db.Close()

	//create a new transaction
	tx := bc.NewUTXOTransaction([]byte(from), []byte(to), amount) 

	bc.MineBlock([]*transactions.Transaction{tx})

	fmt.Println("Success !")

}

func (cli *Cli) createWallet(){
	wallet := wallets.NewWallet()

	fmt.Printf("Your address: %s\n", wallet.GetAddress())
}

//prints the usage of the commands
func (cli *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" printchain - print all the blocks of the blockchain")
	fmt.Println("createchain - creates a chain if none exists")
	fmt.Println("getbalance -address fetches the coins balance for a specific address")
}


//validate that all args are passed
func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}






