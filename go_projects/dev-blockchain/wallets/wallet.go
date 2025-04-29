package wallets

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/elliptic"
	"log"
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




