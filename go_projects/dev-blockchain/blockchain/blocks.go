package blockchain

import (
	"bytes"
	"crypto/sha256"
	"log"
	"time"
	"encoding/gob"
	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/transactions"
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

