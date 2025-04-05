package models

import (
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
	PrivateKey string `json:"privateKey"`
	PublicKeyX string `json:"publicKeyX"`
	PublicKeyY string `json:"publicKeyY"`
	Address    string `json:"address"`
}
type Receiver struct {
	ReceiverPublicKeyX string `json:"privateKey"`
	ReceiverPublicKeyY string `json:"publicKeyX"`
}
type Sender struct {
	SenderPublicKeyY string `json:"privateKey"`
	SenderPublicKeyX string `json:"publicKeyX"`
	SenderPrivateKey string `json:"publicKeyY"`
}

type Transaction struct {
	SenderPrivateKey   string
	SenderPublicKeyX   string
	SenderPublicKeyY   string
	ReceiverPublicKeyX string
	ReceiverPublicKeyY string
	Amount             int
	Signature          string
}
