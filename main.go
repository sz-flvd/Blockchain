package main

import (
	"fmt"
	"krypto.blockchain/src/api"
	"krypto.blockchain/src/threads"
	"krypto.blockchain/src/client"
)

func main() {
	fmt.Println("Jebac krypto")

	// TODO init nodes properly
	nodes := make([]threads.Node, 6)
	for i := 0; i < 6; i++ {
		nodes[i].NewRecordChannel = make(chan []byte, 4)
	}

	manager := api.Manager{Nodes: nodes}
	server := api.BlockchainServer{}

	go client.HandleInput()

	server.Run(&manager)
}

