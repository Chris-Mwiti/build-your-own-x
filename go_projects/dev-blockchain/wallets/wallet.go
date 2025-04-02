package wallets

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"log"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "databases/wallet.dat"

//this is the case scenario of a wallet:
//a wallet contains the following:
type Wallet struct {
	PrivateKey []byte
	PublicKey []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func (w Wallet) SaveToFile(){
	var walletContent bytes.Buffer

	//create a new encoder to encode the data
	encoder := gob.NewEncoder(&walletContent)

	err := encoder.Encode(w)

	if err != nil {
		log.Panicf("error encoding wallet content: %#v", err)
	}

	//write the bytes wallet content to a file
	err = os.WriteFile(walletFile, walletContent.Bytes(), 0644)

	if err != nil {
		log.Panicf("error writing wallet to file: %v", err)
	}
}

//step 1: create a new key pair of keys(private, public)
//creates a new keypair(private key, public key)
//public keys are a point inside the curve
func newKeyPair()([]byte, []byte){
	curve := elliptic.P256()
	private, x, y, err := elliptic.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pubKey := append(x.Bytes(), y.Bytes()...)

	return private, pubKey
}

//creation of a new wallet 
func NewWallet() (*Wallet, error){

	//validate if the wallet already exists
	if _, err := os.Stat(walletFile); !os.IsNotExist(err){
		return nil, errors.New("wallet already exists")
	}

	private, public := newKeyPair()
	wallet := Wallet{
		PrivateKey: private,
		PublicKey: public,
	}

	return &wallet, nil
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

func (wallet Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(wallet.PublicKey)

	versiondedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versiondedPayload)

	fullPayload := append(versiondedPayload, checksum...)
	address := Base58Encode(fullPayload)


	return address
} 




