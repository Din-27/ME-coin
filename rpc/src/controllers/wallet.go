package controllers

import (
	"net/http"
	"github.com/Din-27/blockchain/src/services"
	"github.com/gin-gonic/gin"
)

func HandleCreateWallet(c *gin.Context) {
	wallet := services.NewWallet()
	c.JSON(http.StatusOK, wallet)
}

func HandleCreateWalletSender(c *gin.Context) {
	wallet := services.NewWallet()
	c.JSON(http.StatusOK, wallet)
}
