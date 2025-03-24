package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"github.com/boltdb/bolt"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
)

type Blockchain struct {
	//holds the current hash of the block in the chain
	Tip []byte

	//store the Db connection...used to maintain a connection while the program is running
	Db *bolt.DB
}

//blockchain iterator type
type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}


const dbFile = "databases/blockchain.db"

//holds the key value pairs of the blocks
const blocksBucket = "blocksBucket"


//utility func to check if a db exists
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

//generally the func gives you back the tip of the blockchain
func NewBlockChain(address string) *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db,err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{Tip: tip, Db: db}

	return &bc
}


//creation of a blockchain with db
func CreateBlockchain(address string) *Blockchain {

	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}
	//set the Tip pointer of the current block
	var tip []byte
	db,err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		//create the genesis block from a coinbase transaction
		coinbaseTx := transactions.NewCoinbaseTX(address, transactions.GenesisCoinbaseData)
		genesis := NewGenesisBlock(coinbaseTx);

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}


		//set the key as the genesis hash and the value as the serialized block version
		err = b.Put(genesis.Hash, genesis.Serialze())
		if err != nil {
			return err
		}

		//store the pointer hash key for the block
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		//sets the pointer of the current hash block
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//create a db connected blockchain with the 
	//current and latest block hash & ongoing the db connection
	bc := Blockchain{
		Tip: tip,	
		Db: db,	
	}

	return &bc
}

//Add a new block to the chain
func (bc *Blockchain) MineBlock(transactions []*transactions.Transaction){
	//get the last block added in the chain
	var lastHash []byte
	
	err := bc.Db.View(func(tx *bolt.Tx) error{
		//fetch the blocks buckte and the last block in the chain
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//create a new block with the fetched lastHash
	newBlock := NewBlock(transactions, lastHash)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialze())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transactions.Transaction {
	var unspent []transactions.Transaction

	//stores the spent transactions within a transaction 	
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator();

	for {
		block, err := bci.Next()

		if err != nil {
			log.Panic(err)
			break
		}

		//we have reached the end of the blockchain
		if len(block.PrevBlockHash) == 0 {
			break
		}
		//loop through the transaction in each block
		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Vout {
					//we check whether our local dict of spent transactions contain the transaction
					//if not then we loop over the spent transactions outputs indexex and compare if it matches with out current output index we are
					//if its true we continue to the next transaction index since our focus is finding the unspent transactions  
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}

					if out.IsLockedWithKey(pubKeyHash) {
						unspent = append(unspent, *tx)
					}

				}
				//checks whether the transaction is a coinbase(initial) transaction
				//if not the transaction must have inputs..so we loop over the input
				//and check whether we can unlock the input with the address(but eventually will change)
				//if we unlock we append to the spent transactions slice of that transaction the index of the output being referenced
				if !tx.IsCoinbase() {
					for _, in := range tx.Vin {
						if in.UsesKey(pubKeyHash) {
							inTxId := hex.EncodeToString(in.Txid)
							spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
						}
					}

				}
		}
	}

	return unspent
}

func (bc *Blockchain) FindUnspentTxo(pubKeyHash []byte) []transactions.TxOutput{
	var UTXOS []transactions.TxOutput

	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash){
				UTXOS = append(UTXOS, out)
			}
		}
	}

	return UTXOS
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int)(int, map[string][]int){
	//stores the unspent outputs
	spendableOutputs := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

	//here we loop over all collected unspent transactions
	//and check whether their outputs can be unlocked with current address
	//and the accumulated value is less than the amount checked against
	Work:
		for _, tx := range unspentTxs {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Vout {
				if out.CanBeUnlockedWith(address) && accumulated < amount {
					accumulated += out.Value
					spendableOutputs[txID] = append(spendableOutputs[txID], outIdx)
				}			

				if accumulated >= amount {
					break Work
				}
			}
		}
	
	return accumulated, spendableOutputs
}


//sending coins...here we create a new transaction, put it in a block
//mine the block
func (bc *Blockchain) NewUTXOTransaction(from,to string, amount int) *transactions.Transaction {
    //stores the inputs of the transaction from
    var inputs []transactions.TxInput
    //stores the outputs after a transaction is made
    var outputs []transactions.TxOutput

    //checks and extracts the spendable outputs in a transaction
    acc, validOutputs := bc.FindSpendableOutputs(from, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }


    for txid, outs := range validOutputs {
        txId, err := hex.DecodeString(txid)

        if err != nil {
            log.Panicf("Error while decoding tx id:%v", err)
        }
        for _, out := range outs {
            //creates an input from the outputs from the sender
            input := transactions.TxInput{
                Txid: txId,
                Vout: out,
                Signature: []byte(from),
            }
            inputs = append(inputs, input)
        }
    }

    //Build a list of outputs
    outputs = append(outputs, transactions.TxOutput{
		Value: amount,
		PubKeyHash: []byte(to),
	})

    if acc > amount {
		//we create a change incase the amount exceeds the cumulated amount
        outputs = append(outputs, transactions.TxOutput{
			Value: acc - amount,
			PubKeyHash: []byte(from),
		})
    }

    //create a new transaction based on the generated outputs and inputs
    tx := transactions.Transaction{
		ID: nil,
		Vin: inputs,
		Vout: outputs,
	}

    //create the ID of the transaction
    tx.SetID()

    return &tx
}


//creates an iterator which can be used to traverse through the blocks in the chain
func (bc *Blockchain) Iterator() *BlockchainIterator {
	//the Tip of a blockchain...from the top to the bottom...newest to the oldest
	bci := &BlockchainIterator{
		currentHash: bc.Tip,
		db: bc.Db,
	}

	return bci
}

//@todo: Research on when do we know we have reached the final block in the chain
func (i *BlockchainIterator) Next() (*Block, error) {
	var block *Block
	
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//perfoms a get operation for the current block in the chain
		//deserialize the block from the bytes array to block struct
		encodedblock := b.Get(i.currentHash)
		block = DeserialzeBlock(encodedblock)

		return nil
	})

	log.Printf("Block %#v", block)

	if err != nil {
		return nil,err
	}

	//set the iterator current Hash block pointer..
	//to the prevBlock in the chain
	//we have done this since the latest block is the one added first
	i.currentHash = block.PrevBlockHash
	return block,nil
}