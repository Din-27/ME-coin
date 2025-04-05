package controllers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Din-27/blockchain/src/models"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
)

func padTo32(b []byte) []byte {
	if len(b) < 32 {
		pad := make([]byte, 32-len(b))
		return append(pad, b...)
	}
	return b
}

func decodePrivateKey(encodedPrivateKey string) (*ecdsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	// Buat PrivateKey
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = new(big.Int).SetBytes(decoded)

	// Generate PublicKey dari PrivateKey
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(decoded)

	return privKey, nil
}

func SignTransaction(privateKey *ecdsa.PrivateKey, txHash []byte) string {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash)
	if err != nil {
		log.Fatal(err)
	}

	signature := append(padTo32(r.Bytes()), padTo32(s.Bytes())...)
	return hex.EncodeToString(signature)
}

func verifyTransactionSignature(pubKeyXHex, pubKeyYHex string, txHash []byte, signature string) bool {
	// Decode signature
	sig, err := hex.DecodeString(signature)
	if err != nil {
		log.Fatal("Failed to decode signature:", err)
	}

	if len(sig) != 64 {
		log.Println("🚫 Signature length not 64 bytes, got:", len(sig))
		return false
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])

	x := new(big.Int)
	y := new(big.Int)
	_, ok1 := x.SetString(pubKeyXHex, 16)
	_, ok2 := y.SetString(pubKeyYHex, 16)

	if !ok1 || !ok2 {
		log.Println("🚫 Failed to parse public key X or Y")
		return false
	}

	pubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	fmt.Println("🔍 DEBUG: Verifying signature")
	fmt.Println("txHash:", hex.EncodeToString(txHash))
	fmt.Println("r:", r.String())
	fmt.Println("s:", s.String())
	fmt.Println("Public Key X:", pubKey.X.String())
	fmt.Println("Public Key Y:", pubKey.Y.String())

	result := ecdsa.Verify(&pubKey, txHash, r, s)
	if !result {
		fmt.Println("❌ Signature verification failed")
	} else {
		fmt.Println("✅ Signature verified!")
	}
	return result
}

func HandleSenderAmount(c *gin.Context) {
	var tx models.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	txHash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%d", tx.SenderPublicKeyX, tx.SenderPublicKeyY, tx.Amount)))
	decodedPrivateKey, err := decodePrivateKey(tx.SenderPrivateKey)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid sender private key"})
		return
	}

	signature := SignTransaction(decodedPrivateKey, txHash[:])
	tx.Signature = signature

	if verifyTransactionSignature(tx.SenderPublicKeyX, tx.SenderPublicKeyY, txHash[:], signature) {
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
