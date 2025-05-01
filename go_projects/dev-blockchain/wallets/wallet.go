package wallets

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "databases/wallet.dat"

//this is the case scenario of a wallet:
//a wallet contains the following:
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte 
}

//struct to represent the encoded private key of a wallet
type _PrivateKey struct {
	D *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
}

//step 1: create a new key pair of keys(private, public)
//creates a new keypair(private key, public key)
//public keys are a point inside the curve
func newKeyPair()(ecdsa.PrivateKey, []byte){
	curve := elliptic.P256()
	private,err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)


	return *private, pubKey
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


func HashPubKey(pubkey []byte) []byte{
	publicSHA256 := sha256.Sum256(pubkey)	

	RIPEMD160Hasher := ripemd160.New()
	_,err := RIPEMD160Hasher.Write(publicSHA256[:])

	if err != nil {
		log.Panicf("Error while hashing the public key: %v", err)
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
	//hash the public key
	pubKeyHash := HashPubKey(wallet.PublicKey)

	//append the version payload as a prefix
	versiondedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versiondedPayload)

	fullPayload := append(versiondedPayload, checksum...)
	address := Base58Encode(fullPayload)


	return address
} 

//a work around the error while encoding the wallets

func (wallet *Wallet) GobEncode()([]byte, error){
	privKey := &_PrivateKey{
		D: wallet.PrivateKey.D,
		PublicKeyX: wallet.PrivateKey.X,
		PublicKeyY: wallet.PrivateKey.Y,
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(privKey)

	if err != nil {
		return nil, err
	}

	_, err = buf.Write(wallet.PublicKey)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//decodes the wallet 
func (wallet *Wallet) GobDecode(data []byte)(error){
	buf := bytes.NewBuffer(data)

	var privKey _PrivateKey
	
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&privKey)
	if err != nil {
		return err
	}

	wallet.PrivateKey = ecdsa.PrivateKey{
		D: privKey.D,
		PublicKey: ecdsa.PublicKey{
			X: privKey.PublicKeyX,
			Y: privKey.PublicKeyY,
			Curve: elliptic.P256(),
		},
	}

	wallet.PublicKey = make([]byte, buf.Len())
	_, err = buf.Read(wallet.PublicKey)
	if err != nil {
		return err
	}

	return nil
}




