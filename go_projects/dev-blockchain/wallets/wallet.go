package wallets

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "databases/wallet.dat"

//this is the case scenario of a wallet:
//a wallet contains the following:
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey ecdsa.PublicKey
}

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

//step 1: create a new key pair of keys(private, public)
//creates a new keypair(private key, public key)
//public keys are a point inside the curve
func newKeyPair()(ecdsa.PrivateKey, ecdsa.PublicKey){
	curve := elliptic.P256()
	private,err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}


	return *private, private.PublicKey
}

//creation of a new wallet 
func NewWallet() (*Wallet){
	private, public := newKeyPair()
	wallet := Wallet{
		PrivateKey: private,
		PublicKey: public,
	}

	return &wallet
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

func HashPubKey(pubkey []byte) []byte{
	publicSHA256 := sha256.Sum256(pubkey)	

	RIPEMD160Hasher := ripemd160.New()
	_,err := RIPEMD160Hasher.Write(publicSHA256[:])

	if err != nil {
		log.Panicf("Error[HashPubKey]: %v", err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]

}

//creates a human readable address for the wallet public address
func (wallet Wallet) GetAddress() []byte {

	//wallet address public key 
	publicKey := append(wallet.PublicKey.X.Bytes(), wallet.PublicKey.Y.Bytes()...)

	//hash the public key
	pubKeyHash := HashPubKey(publicKey)

	//append the version payload as a prefix
	versiondedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versiondedPayload)

	fullPayload := append(versiondedPayload, checksum...)
	address := Base58Encode(fullPayload)


	return address
} 




