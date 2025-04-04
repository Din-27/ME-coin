package main

import (
	"fmt"
	"github.com/Din-27/blockchain/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.Routes(router)

	fmt.Println("Blockchain RPC Server running on :8080")
	router.Run()
}
