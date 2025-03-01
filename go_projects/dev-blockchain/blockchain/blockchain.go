package blockchain

import (
	"bytes"
	"crypto/sha256"
	"log"
	"time"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
	"github.com/boltdb/bolt"
)

//block -> store valuable information...
type Block struct {
	PrevBlockHash []byte
	Hash []byte
	Timestamp int64
	Transaction []*transactions.Transaction
	Nounce int
}


//Tx struct that will keep track of all the transactions
type Tx struct {
	Data []byte
}

type Blockchain struct {
	//holds the current hash of the block in the chain
	Tip []byte

	//store the Db connection...used to maintain a connection while the program is running
	Db *bolt.DB
}


//block creation
func NewBlock(transactions []*transactions.Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Transaction: transactions,
		PrevBlockHash: prevBlockHash,
		Hash: []byte{},
		Nounce:0,
	}

	pow := NewProofOfWork(block)
	//run to capture the hash and set it to the headers
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nounce = nonce

	return block
}

//Add a new block to the chain
func (bc *Blockchain) AddBlock(data string){
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
	newBlock := NewBlock(data, lastHash)

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

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash  [32]byte

	//store each id in the transaction in to one hash
	for _, tx := range block.Transaction {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
	
}


//func that actually creates the genesis block
func NewGenesisBlock(coinbase *transactions.Transaction) *Block {
	return NewBlock([]*transactions.Transaction{
		coinbase,
	}, []byte{})
}
//creates a new blockchain with the actual blockchain
func NewBlockchain() *Blockchain {
	return BlockChainWithDb(transactions.GenesisCoinbaseData)
}

