package wallets

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string][]byte
}

//encrypts the content of the wallets and saves it in a .dat file format
func (ws Wallets) SaveToFile(){
	var walletContent bytes.Buffer

	//create a new encoder to encode the data
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
		return err
	}

	var wallet Wallets
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallet)

	if err != nil {
		return err
	}

	ws.Wallets = wallet.Wallets

	return nil

}
//returns the already created wallets and loads them as a list of wallets
func WalletsList() (*Wallets, error) {
	wallets := Wallets{}

	wallets.Wallets= make(map[string][]byte)

	err := wallets.LoadFromFile();

	return &wallets, err 
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

	//encode the wallet for storage
	encWallet, err := wallet.GobEncode()
	
	if err != nil {
		log.Panicf("Error while encoding wallet: %v", err)
	}


	//attaches the newly created wallet to the wallets collection 
	ws.Wallets[address] = encWallet 

	return address
}

func (ws Wallets) GetWallet(address  string) Wallet {
	var wallet Wallet
	wallet.GobDecode(ws.Wallets[address])
	return wallet
}
