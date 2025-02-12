package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//nodes taking part in the blockchain
type Account string


//rep of the data struct of a transaction
type Tx struct {
	From Account `json:"from"`
	To Account `json:"to"`
	Value uint `json:"value"`
	Data string `json:"data"`
}

//func that will indicate whether a transaction is a reward or not
func (tx Tx) IsReward()bool {
	return tx.Data == "reward"
}

/** 
@todo: remember to create a func that will parse the geneseis transaction
@todo: remember to register each transaction as an event
*/

//DB component that will encapsulate all db state
//know all about all users balances and who transferred TBB tokens to whom
//1. Adding new transactions to Mempool
//2. Validate transactions against current state
//3. Changing state
//4. Persisiting transactions to disk
//5. Calculating accounts balances by replaying all transactions since genesis in a sequence
type State struct {
	Balances map[Account]uint
	txMempool []Tx

	//event file that will keep track of all transactins
	dbFile *os.File
}

func NewStateFromDisk() (*State, error) {
	//get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		//update the balances 
		balances[account] = balance
	}

	//genesis state balances are updated by sequentially replaying all the database events
	txDbFilePath := filepath.Join(cwd, "databases", "txt.db")

	//generate the fd 
	fd, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	//create a new scanner that will read all the transaction events
	scanner := bufio.NewScanner(fd)
	state := &State{
		Balances: balances,
		txMempool: make([]Tx, 0),
		dbFile: fd,
	}

	//iterate over each the tx.db files line
	for scanner.Scan(){
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	//convert JSON encoded TX into an object (struct)
	var tx Tx
	json.Unmarshal(scanner.Bytes(), &tx)


	//Rebuild the state (user balances), 
	//as a series of events
	if err := state.apply(tx); err != nil {
		return nil, err
	}
	
	return state, nil
}

//Adding a new transaction to th Mempool
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil  {
		return err
	}

	s.txMempool = append(s.txMempool, tx)
	return nil
}

//Persisting the transactions to disk
func (state *State) Persist() error {
	//make a copy of mempool because s.txMempool will be modified
	
	//creation of a mempool
	mempool := make([]Tx, len(state.txMempool))
	// create a copy of the mempool and save it
	copy(mempool, state.txMempool)

	for i:= 0; i < len(state.txMempool); i++{
		txJson, err := json.Marshal(mempool[i]) 
		if err != nil {
			return err
		}

		//persist the jsonFormData to a the dbfile  
		if _, err = state.dbFile.Write(append(txJson, '\n')); err != nil {
			return err	
		}

		//Remove the TX written to a file from the mempool
		state.txMempool = state.txMempool[1:]
	}

	return nil
}

func (state *State) apply(tx Tx) error {
	//check if its a reward
	if tx.IsReward() {
		state.Balances[tx.To] += tx.Value
	}

	//
	if tx.Value > state.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	//decrementation of acc bal from the sender
	state.Balances[tx.From] -= tx.Value
	state.Balances[tx.To] += tx.Value


	return nil
} 


