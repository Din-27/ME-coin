package controllers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/Din-27/blockchain/src/services"
	"github.com/gin-gonic/gin"
)

func decodePrivateKey(encodedPrivateKey string) (*ecdsa.PrivateKey, error) {
	// Decode private key dari base64
	decoded, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	// Membuat private key dari byte array decoded
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKey.D = new(big.Int).SetBytes(decoded[:len(decoded)-1])
	privKey.PublicKey.X = new(big.Int).SetBytes(decoded[len(decoded)-1:])
	privKey.PublicKey.Y = new(big.Int).SetBytes(decoded[len(decoded)-1:])
	privKey.PublicKey.Curve = elliptic.P256()

	return privKey, nil
}

func SignTransaction(privateKey *ecdsa.PrivateKey, txHash []byte) string {
	// Menandatangani transaksi menggunakan private key
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash)
	if err != nil {
		log.Fatal(err)
	}

	// Menggabungkan r dan s menjadi signature yang akan dikirim ke client
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature)
}

func verifyTransactionSignature(publicKey []byte, txHash []byte, signature string) bool {
	// Decode signature dari hex string
	sig, err := hex.DecodeString(signature)
	if err != nil {
		log.Fatal(err)
	}

	// Membagi signature menjadi dua bagian: r dan s
	r := new(big.Int).SetBytes(sig[:len(sig)/2])
	s := new(big.Int).SetBytes(sig[len(sig)/2:])

	// Membagi public key menjadi bagian X dan Y
	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(publicKey[:len(publicKey)/2]),
		Y:     new(big.Int).SetBytes(publicKey[len(publicKey)/2:]),
	}

	// Memverifikasi signature
	return ecdsa.Verify(&pubKey, txHash, r, s)
}

func HandleSenderAmount(c *gin.Context) {
	var tx services.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Membuat hash untuk transaksi
	txHash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%d", tx.Sender, tx.Receiver, tx.Amount)))

	// Mendekode public key pengirim
	senderPublicKey, err := hex.DecodeString(tx.Sender)
	fmt.Println(tx.Sender)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid sender public key"})
		return
	}

	// Mendekode private key pengirim dari body
	decodedPrivateKey, err := decodePrivateKey(tx.SenderPrivateKey)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid sender private key"})
		return
	}

	// Menandatangani transaksi dengan private key pengirim
	signature := SignTransaction(decodedPrivateKey, txHash[:])

	// Menambahkan signature ke dalam transaksi
	tx.Signature = signature

	// Memverifikasi transaksi
	if verifyTransactionSignature(senderPublicKey, txHash[:], signature) {
		c.JSON(200, gin.H{
			"status":      "Transaction verified successfully!",
			"transaction": tx,
		})
	} else {
		c.JSON(400, gin.H{
			"status": "Transaction verification failed!",
		})
	}
}
