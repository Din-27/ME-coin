package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Din-27/blockchain/src/models"
	"log"
)

func NewWallet() *models.Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	xBytes := privateKey.PublicKey.X.Bytes()
	yBytes := privateKey.PublicKey.Y.Bytes()

	// Buat address dari gabungan X + Y public key
	publicKey := append(xBytes, yBytes...)
	address := generateAddress(publicKey)

	return &models.Wallet{
		PrivateKey: hex.EncodeToString(privateKey.D.Bytes()),
		PublicKeyX: hex.EncodeToString(xBytes),
		PublicKeyY: hex.EncodeToString(yBytes),
		Address:    address,
	}
}

func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[:])
}
