package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

//block -> store valuable information...
type Block struct {
	PrevBlockHash []byte
	Hash []byte
	Timestamp int64
	Data *Tx 
}


//Tx struct that will keep track of all the transactions
type Tx struct {
	Data []byte
}

type Blockchain struct {
	blocks []*Block
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
	}

	block.SetHash()
	return block
}

//Add a new block to the chain
func (bc *Blockchain) AddBlock(data string){
	//get the last block added in the chain
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}


//func that actually creates the genesis block
func NewGenesisBlock() *Block {
	return NewBlock("The genesis block", []byte{})
}
//creates a new blockchain with the actual blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main(){
	chain := NewBlockchain()

	chain.AddBlock("First member to join the chain")
	chain.AddBlock("Second member to join the chain")
	
	for _, block := range chain.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)	
		fmt.Printf("Data: %s\n", block.Data.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}