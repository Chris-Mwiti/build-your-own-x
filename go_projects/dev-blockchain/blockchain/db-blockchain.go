package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
	"github.com/boltdb/bolt"
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


const dbFile = "databases/blocks.db"

//holds the key value pairs of the blocks
const blocksBucket = "blocksBucket"


//utility func to check if a db exists
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}


//creation of a blockchain with db
func BlockChainWithDb(address string) *Blockchain {

	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	//set the Tip pointer of the current block
	var tip []byte
	db,err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		//check if the blocks bucket already exists
		if b == nil {
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
		} else {
			//fetch the last hash block instance in the chain
			tip = b.Get([]byte("l"))
		}

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

func (bc *Blockchain) FindUnspentTransactions(address string) []transactions.Transaction {
	var unspent []transactions.Transaction

	//research more on this how does it store data
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator();

	for {
		block, err := bci.Next()

		if err != nil {
			log.Panic(err)
			break
		}
		//loop through the transaction in each block
		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Vout {
					//was the output spent?
					if spentTXOs[txID] != nil {

					}
				}
		}
	}
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
	
	err := i.db.View(func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//perfoms a get operation for the current block in the chain
		//deserialize the block from the bytes array to block struct
		encodedblock := b.Get(i.currentHash)
		block = DeserialzeBlock(encodedblock)

		return nil
	})

	if err != nil {
		return nil,err
	}

	//set the iterator current Hash block pointer..
	//to the prevBlock in the chain
	//we have done this since the latest block is the one added latest
	i.currentHash = block.PrevBlockHash
	return block,nil
}