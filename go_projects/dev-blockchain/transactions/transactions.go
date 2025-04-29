package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"github.com/Chris-Mwiti/build-your-own-x/go-projects/dev-blockchain/wallets"
)

//bitcoin transactions do not store the following:
//1. No accounts
//2. No balances
//3. No addresses
//4. No coins
//5. No senders and receivers

//example of a bitcoin structure type
type Transaction struct {
    ID []byte
    Vin []TxInput //references of transactions that are inputs to this transactions
    Vout []TxOutput //stamps the value 
}

//input of a new trasaction reference the outputs of a previous transaction
//transactions just lock values with a script which can be unlocked only by the one who locked them

type TxOutput struct {
    //actually stores "coins"
    Value int
    //locks the transaction with a puzzle
    PubKeyHash []byte //will store user defined wallet addresses for now
}

type TxInput struct {
    Txid []byte  //store the id of the transaction being referenced
    Vout int //stores an index of an output in the transaction
    Signature []byte //provides data to be used in the ScriptPubKey ...if data is correct, the output can be unlocked, and its value can be used to generate new outputs 
    PubKey []byte
}

//genesis block data
const GenesisCoinbaseData = "17th February 2025"

func (tx *Transaction) SetID(){
    var encoded bytes.Buffer;
    var hash [32]byte

    enc := gob.NewEncoder(&encoded);
    err := enc.Encode(tx);

    if err != nil {
        log.Panic(err)
    }

    hash = sha256.Sum256(encoded.Bytes())

    tx.ID = hash[:]
}

//coinbase transaction
//special type of transaction which doesn't require previously existing outputs
func NewCoinbaseTX(to, data string) *Transaction {
    if data == ""{
        data = fmt.Sprintf("Reward to %s", to)
    }
    
    txin := TxInput{
        Txid: []byte{},
        Vout: -1,
        Signature: []byte(data),
    }

    //@todo: implement a proper subsidy strategy
    subsidy := 20;
    txout := TxOutput{
        Value: subsidy,
        PubKeyHash: []byte(to),
    }

    tx := Transaction{
        ID: nil,
        Vin: []TxInput{txin},
        Vout: []TxOutput{txout},
    }
    

    //set the id for the transaction
    tx.SetID()

    fmt.Printf("New transactions: %#v", tx.Vout)

    return &tx
}


//unspent transactions section
//unspent transactions means that these outputs weren't referenced in any inputs
//we can only unlock those that can be unlocked by the key we own


//checks if a transaction is a coinbase
func (tx *Transaction) IsCoinbase() bool {
    if len(tx.Vin) != 1 {
        return false
    }
    return tx.Vin[0].Vout == -1
}


//stringify the transaction to make it human readable
func (tx Transaction) String() (string) {
    var lines []string


    //prepend the transaction id to the begining of line
    lines = append(lines, fmt.Sprintf("Transaction Id: %x",tx.ID ))

    //here we create information txt structure for the inputs
    for i, input := range tx.Vin {
        lines = append(lines, fmt.Sprintf("Input: %d", i))
        lines = append(lines, fmt.Sprintf("TxInputId: %x", input.Txid))
        lines = append(lines, fmt.Sprintf("Output reference: %d", input.Vout))
        lines = append(lines, fmt.Sprintf("ScriptSig: %s", input.Signature))
    }

    for i, outputs := range tx.Vout {
        lines = append(lines, fmt.Sprintf("Output: %d", i))
        lines = append(lines, fmt.Sprintf("Value: %d", outputs.Value))
        lines = append(lines, fmt.Sprintf("Script: %s", outputs.PubKeyHash))
    }

    return strings.Join(lines, "\n")
}



//an input has a pubkey field to identify the receiver of the transaction
//and a signature to verify the sender of the transaction
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
    lockingHash := wallets.HashPubKey(in.PubKey)

    return bytes.Equal(lockingHash, pubKeyHash) 
}

//locks the output of a transaction with memonic base58 wallet address
func (out *TxOutput) Lock(address []byte) {
    pubKeyHash := wallets.Base58Encode(address)
    //slices of the checksum of the transaction and only returns (version & pubKeyHash)
    pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
    out.PubKeyHash = pubKeyHash
}


//verify's whether the output of a transaction has been locked with the specified pubKey
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
    return bytes.Equal(out.PubKeyHash, pubKeyHash) 
}