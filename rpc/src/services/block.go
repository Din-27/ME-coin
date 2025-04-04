package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Din-27/blockchain/src/models"
	"math/big"
	"os"
	"strconv"
	"time"
)

var difficulty = 2 // Ubah untuk menyesuaikan kesulitan mining

func CalculateHash(block models.Block) string {
	blockData := strconv.Itoa(block.Index) + block.Timestamp + block.Transactions + block.PrevHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

func MineBlock(prevBlock models.Block, transactions string) models.Block {
	var newBlock models.Block
	newBlock.Index = prevBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = prevBlock.Hash

	for {
		newBlock.Nonce++
		newBlock.Hash = CalculateHash(newBlock)
		if IsValidHash(newBlock.Hash) {
			break
		}
	}
	return newBlock
}

func IsValidHash(hash string) bool {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty*4))
	hashInt := new(big.Int)
	hashInt.SetString(hash, 16)
	return hashInt.Cmp(target) == -1
}

func SaveBlockchain(filename string, blockchain models.Blockchain) {
	file, _ := json.MarshalIndent(blockchain, "", "  ")
	_ = os.WriteFile(filename, file, 0644)
}

func LoadBlockchain(filename string) models.Blockchain {
	var bc models.Blockchain
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Blockchain file not found, creating new one.")
		return models.Blockchain{Blocks: []models.Block{{0, time.Now().String(), "Genesis Block", "", 0, ""}}}
	}
	_ = json.Unmarshal(file, &bc)
	return bc
}
