package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"github.com/boltdb/bolt"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
)

//blockchain iterator type
type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

const dbFile = "databases/blocks.db"

//holds the key value pairs of the blocks
const blocksBucket = "blocksBucket"

//serialization of the block into a byte array a format that can be stored
//in the boltdb
func (b *Block) Serialze() []byte {
	//create a new buffer that will store the bytes array
	var result bytes.Buffer

	//create a new encoder that will encode the data into byte array
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	//check if the transmitted is a nil pointer
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//deserialize func that will revert the byte array into a *Block struct
//this will be an independent function
func DeserialzeBlock(d []byte) *Block {
	var block Block

	//init a new decoder 
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)

	return &block
}

//this function is used to create the first block in the blockchain
func CreateBlockchain(address string) *Block {
	//create a coinbase contex of the genesis block
	cbtx := transactions.NewCoinbaseTX(address, transactions.GenesisCoinbaseData);
	
	//creation of the first genesis block of the chain
	genesis := NewGenesisBlock(cbtx);

	return genesis
}

func BlockChainWithDb(address string) *Blockchain {
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
			//create the genesis block
			genesis := CreateBlockchain(address)
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