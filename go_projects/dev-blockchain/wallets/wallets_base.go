package wallets

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"crypto/elliptic"
	"fmt"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

//encrypts the content of the wallets and saves it in a .dat file format
func (ws Wallets) SaveToFile(){
	var walletContent bytes.Buffer

	//create a new encoder to encode the data
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&walletContent)

	err := encoder.Encode(ws)

	if err != nil {
		log.Panicf("Error encoding wallet: %v", err)
	}

	//write the bytes wallet content to a file
	err = os.WriteFile(walletFile, walletContent.Bytes(), 0644)

	if err != nil {
		log.Panicf("error writing wallet to file: %v", err)
	}
}

//decode the wallets content and attach it to the current wallets object
func (ws *Wallets) LoadFromFile() (error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err){
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panicf("Error while loading wallet file: %v", err)
	}

	var wallet Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallet)

	if err != nil {
		log.Panicf("Error while decoding wallet file: %v", err)
	}

	ws.Wallets = wallet.Wallets

	return nil

}
//returns the already created wallets and loads them as a list of wallets
func WalletsList() *Wallets {
	wallets := Wallets{}

	wallets.Wallets= make(map[string]*Wallet)

	err := wallets.LoadFromFile();
	if err != nil {
		log.Panic("Wallets file doesn't exist")
	}

	return &wallets
}

func (ws *Wallets) ListAddress() ([]string){
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

//creates a wallet and attaches it to the already existing wallets
func (ws *Wallets) CreateWallet() string {
	//create a new wallet with the public and private key and obtain the address
	wallet := NewWallet()

	//load wallets stored content
	address := fmt.Sprintf("%s", wallet.GetAddress())


	//attaches the newly created wallet to the wallets collection 
	ws.Wallets[address] = wallet

	return address
}
