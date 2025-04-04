package routes

import (
	"github.com/Din-27/blockchain/src/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.POST("/mine", controllers.HandleMineBlock)
	router.GET("/blockchain", controllers.HandleGetBlockchain)
	router.GET("/generate-wallet", controllers.HandleCreateWallet)
	router.GET("/generate-wallet-sender", controllers.HandleCreateWalletSender)
	router.POST("/sender", controllers.HandleSenderAmount)
}
