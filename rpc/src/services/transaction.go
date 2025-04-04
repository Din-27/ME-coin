package services

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Din-27/blockchain/src/models"
)

type Transaction struct {
	SenderPrivateKey    string
	Sender    string
	Receiver  string
	Amount    int
	Signature string
}

func NewTransaction(senderWallet *models.Wallet, receiver string, amount int) *Transaction {
	transaction := &Transaction{
		Sender:   senderWallet.Address,
		Receiver: receiver,
		Amount:   amount,
	}
	transaction.SignTransaction(senderWallet.PrivateKey)
	return transaction
}

func (t *Transaction) SignTransaction(privateKey *ecdsa.PrivateKey) {
	hash := sha256.Sum256([]byte(t.Sender + t.Receiver + fmt.Sprintf("%d", t.Amount)))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Fatal(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	t.Signature = hex.EncodeToString(signature)
}

func SignTransaction(privateKey *ecdsa.PrivateKey, txHash []byte) string {
	// Menandatangani transaksi dengan private key
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash)
	if err != nil {
		log.Fatal(err)
	}

	// Menggabungkan r dan s menjadi signature
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature)
}