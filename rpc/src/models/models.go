package models

import (
	"crypto/ecdsa"
	"sync"
)

type Block struct {
	Index        int    `json:"index"`
	Timestamp    string `json:"timestamp"`
	Transactions string `json:"transactions"`
	PrevHash     string `json:"prev_hash"`
	Nonce        int    `json:"nonce"`
	Hash         string `json:"hash"`
}

type Blockchain struct {
	Blocks []Block `json:"blocks"`
	Mux    sync.Mutex
}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	Signature string
}
