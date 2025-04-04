package controllers

import (
	"net/http"
	"github.com/Din-27/blockchain/src/services"
	"github.com/gin-gonic/gin"
)

var blockchain = services.LoadBlockchain("json/blockchain.json")

func HandleGetBlockchain(c *gin.Context) {
	blockchain.Mux.Lock()
	defer blockchain.Mux.Unlock()
	c.JSON(http.StatusOK, blockchain)
}

func HandleMineBlock(c *gin.Context) {
	blockchain.Mux.Lock()
	defer blockchain.Mux.Unlock()

	newBlock := services.MineBlock(blockchain.Blocks[len(blockchain.Blocks)-1], "Miner Reward")
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
	services.SaveBlockchain("src/json/blockchain.json", blockchain)

	c.JSON(http.StatusOK, newBlock)
}
