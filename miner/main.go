package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
)

type Block struct {
	Index        int    `json:"index"`
	Timestamp    string `json:"timestamp"`
	Transactions string `json:"transactions"`
	PrevHash     string `json:"prev_hash"`
	Nonce        int    `json:"nonce"`
	Hash         string `json:"hash"`
}

var difficulty = 2 // Harus sama dengan yang ada di server

func calculateHash(block Block) string {
	blockData := strconv.Itoa(block.Index) + block.Timestamp + block.Transactions + block.PrevHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}

func isValidHash(hash string) bool {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty*4))
	hashInt := new(big.Int)
	hashInt.SetString(hash, 16)
	return hashInt.Cmp(target) == -1
}

func mineBlock(block Block) Block {
	fmt.Println("Starting mining...")
	for {
		block.Nonce++
		block.Hash = calculateHash(block)
		if isValidHash(block.Hash) {
			fmt.Println("Block mined! Hash:", block.Hash)
			return block
		}
	}
}

func getWork(serverURL string) (Block, error) {
	resp, err := http.Get(serverURL + "/blockchain")
	if err != nil {
		return Block{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var blockchain struct {
		Blocks []Block `json:"blocks"`
	}
	json.Unmarshal(body, &blockchain)

	if len(blockchain.Blocks) == 0 {
		return Block{}, fmt.Errorf("blockchain kosong")
	}
	return blockchain.Blocks[len(blockchain.Blocks)-1], nil
}

func submitBlock(serverURL string, block Block) {
	jsonBlock, _ := json.Marshal(block)
	resp, err := http.Post(serverURL+"/mine", "application/json", bytes.NewReader(jsonBlock))
	if err != nil {
		fmt.Println("Error submitting block:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Block submitted!")
}

func main() {
	serverURL := "http://localhost:8080"
	fmt.Println("Connecting to server", serverURL)

	prevBlock, err := getWork(serverURL)
	if err != nil {
		fmt.Println("Error getting work:", err)
		return
	}

	newBlock := mineBlock(prevBlock)
	submitBlock(serverURL, newBlock)
}
