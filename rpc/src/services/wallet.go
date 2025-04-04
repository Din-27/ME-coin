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

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &models.Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[:])
}
