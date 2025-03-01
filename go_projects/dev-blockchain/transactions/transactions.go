package transactions

import "fmt"

//bitcoin transactions do not store the following:
//1. No accounts
//2. No balances
//3. No addresses
//4. No coins
//5. No senders and receivers

//example of a bitcoin structure type
type Transaction struct {
    ID []byte
    Vin []TxInput
    Vout []TxOutput
}

//input of a new trasaction reference the outputs of a previous transaction
//transactions just lock values with a script which can be unlocked only by the one who locked them

type TxOutput struct {
    //actually stores "coins"
    Value int
    //locks the transaction with a puzzle
    ScriptPubKey string //will store user defined wallet addresses for now
}

type TxInput struct {
    Txid []byte  //store the id of the transaction being referenced
    Vout int //stores an index of an output in the transaction
    ScriptSig string //provides data to be used in the ScriptPubKey ...if data is correct, the output can be unlocked, and its value can be used to generate new outputs 
}

const GenesisCoinbaseData = "This is the first block created in the blockchain"

//coinbase transaction
//special type of transaction which doesn't require previously existing outputs
func NewCoinbaseTX(to, data string) *Transaction {
    if data == ""{
        data = fmt.Sprintf("Reward to %s", to)
    }
    
    txin := TxInput{
        Txid: []byte{},
        Vout: -1,
        ScriptSig: data,
    }

    //@todo: implement a proper subsidy strategy
    subsidy := 0;
    txout := TxOutput{
        Value: subsidy,
        ScriptPubKey: to,
    }

    tx := Transaction{
        ID: nil,
        Vin: []TxInput{txin},
        Vout: []TxOutput{txout},
    }
    

    //set the id for the transaction
    tx.SetID()

    return &tx
}

