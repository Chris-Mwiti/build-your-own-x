package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

//block -> store valuable information...
type Block struct {
	PrevBlockHash []byte
	Hash []byte
	Timestamp int64
	Data *Tx
	Nounce int
}


//Tx struct that will keep track of all the transactions
type Tx struct {
	Data []byte
}

type Blockchain struct {
	//holds the current hash of the block in the chain
	tip []byte

	//store the db connection...used to maintain a connection while the program is running
	db *bolt.DB
}

//creation of hashes of blocks...This is will be used to keep track of blocks
//and make it difficult to actually add block into the network
func (b *Block) SetHash(){
	//creation of a timestamp that will keep track of time a block is created
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//combination of the prevBlockHash,data and timestamp into one byte slice
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data.Data, timestamp}, []byte{})
	//creation of a hash from the headers
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]

}

//block creation
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Data: &Tx{
			Data: []byte(data),
		},
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
	
	err := bc.db.View(func(tx *bolt.Tx) error{
		//fetch the blocks buckte and the last block in the chain
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	//create a new block with the fetched lastHash
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialze())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}


//func that actually creates the genesis block
func NewGenesisBlock() *Block {
	return NewBlock("The genesis block", []byte{})
}
//creates a new blockchain with the actual blockchain
func NewBlockchain() *Blockchain {
	return BlockChainWithDb()
}

func main(){
	chain := NewBlockchain()

	chain.AddBlock("First member to join the chain")
	chain.AddBlock("Second member to join the chain")
	
	//current hash in the block 

	for {

		nextBlck, err := chain.Iterator().Next()

		currBlck := chain.Iterator().currentHash	
		nextBlckHash := nextBlck.Hash
		
		if err != nil {
			break
		}

		fmt.Printf("Current Block: %x\n", currBlck)
		fmt.Printf("Next block: %x\n", nextBlckHash)
	}


}