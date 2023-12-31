package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BlockchainServer struct {
	manager *Manager
}

// Rest api is probably an overkill, but I wanted to have a bit of fun I guess, sorry
// You can use from console with client or via Postman. Base Url: localhost:8080/blockchain
// Example"
//  GET localhost:8080/blockchain/local/8
//  Response { "error": "Incorrect node id: 8" }
func (server *BlockchainServer) Run(blockchainManager *Manager) {
	gin.DisableConsoleColor()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	server.manager = blockchainManager
	router := gin.Default()

	router.GET("/blockchain/:id", server.getBlock)
	router.GET("/blockchain/local/:index", server.getLocalChain)
	router.GET("/blockchain/count", server.getCount)
	router.POST("/blockchain", server.postRecord)
	router.GET("/blockchain/quit", server.quit)
	// Cryptocurrency
	router.GET("/blockchain/public/:idx", server.getPublicKey)

	router.Run("localhost:8080")
}

func (server *BlockchainServer) getBlock(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	response, err := server.manager.GetBlock(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, response)
	}
}

func (server *BlockchainServer) getLocalChain(c *gin.Context) {
	param := c.Param("index")
	index, err := strconv.Atoi(param)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	response, err := server.manager.GetLocalBlockchain(index)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, response)
	}
}

func (server *BlockchainServer) getCount(c *gin.Context) {
	count := server.manager.GetBlockCount()

	c.IndentedJSON(http.StatusOK, count)
}

func (server *BlockchainServer) postRecord(c *gin.Context) {
	var request AddRecordRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := server.manager.AddRecord(request)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, response)
	}
}

func (server *BlockchainServer) quit(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Quitting!")
	go executeQuit()
}

// Cryptocurrency
func (server *BlockchainServer) getPublicKey(c *gin.Context) {
	param := c.Param("idx")
	index, err := strconv.Atoi(param)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	response, err := server.manager.GetPublicKey(index)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// fmt.Println(response.PublicKey.X.String())
		c.IndentedJSON(http.StatusOK, response.PublicKey.X)
	}

	// fmt.Println(string(response))

	// count := server.manager.GetBlockCount()
	// c.IndentedJSON(http.StatusOK, count)
}

// I'm extremely proud of this function
func executeQuit() {
	fmt.Println("Closing all threads...")
	time.Sleep(2 * time.Second)
	fmt.Println("All threads closed, quitting...")
	os.Exit(0)
}
