package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
)

//proof of work algorithim => (HashCash Algo)[https://en.wikipedia.org/wiki/Hashcash]
//steps:
//1. Take some publicly know data(in the case of bitcoin is the block headers)
//2. Add a counter to it the counter starts at 0
//3. Get a hash of the data + counter combination
//4. Check that the hash meets the provided requirements
//increase the counter and repeate the steps 3 and 4 if condition is not met

//block header storing the difficulty at which the block was mined
const targetBits = 24

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//@todo: confirm whether this will actually work
	target = target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{
		block: b,
		target: target,
	}
	
	return pow
}

func IntToHex(num int64) []byte {
	//write the binary representation of the num
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//receiver func that prepares data
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	
	return data
}



