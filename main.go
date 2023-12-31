package main

import (
	"flag"
	"fmt"
	"sync"

	"krypto.blockchain/src/api"
	"krypto.blockchain/src/client"
	"krypto.blockchain/src/common"
	"krypto.blockchain/src/threads"
)

func createNodes(n uint) []*threads.Node {
	nodes := make([]*threads.Node, n)

	blockConfirmChannels, blockMinedChannels, recordAddChannels, recordConfirmChannels := initChannels(n)

	for i := uint(0); i < n; i++ {
		nodes[i] = threads.Node_CreateNode(
			uint(i),
			n,
			blockMinedChannels[i],
			blockConfirmChannels[i],
			recordAddChannels[i],
			recordConfirmChannels[i],
			blockMinedChannels,
			blockConfirmChannels,
			recordAddChannels,
			recordConfirmChannels)
	}

	return nodes
}

func initChannels(n uint) ([]chan threads.MessageBool, []chan threads.MessageBlock, []chan threads.RecordAdd, []chan *common.Record) {
	// these are buffered channel of size n - dunno, thought it might be a good number
	blockConfirmChannels := make([]chan threads.MessageBool, n)
	blockMinedChannels := make([]chan threads.MessageBlock, n)
	recordAddChannels := make([]chan threads.RecordAdd, n)
	recordConfirmChannels := make([]chan *common.Record, n)
	for i := uint(0); i < n; i++ {
		blockConfirmChannel := make(chan threads.MessageBool)
		blockConfirmChannels[i] = blockConfirmChannel

		blockMinedChannel := make(chan threads.MessageBlock)
		blockMinedChannels[i] = blockMinedChannel

		recordAddChannel := make(chan threads.RecordAdd)
		recordAddChannels[i] = recordAddChannel

		recordConfirmChannel := make(chan *common.Record)
		recordConfirmChannels[i] = recordConfirmChannel
	}

	return blockConfirmChannels, blockMinedChannels, recordAddChannels, recordConfirmChannels
}

func initThreads(nodes []*threads.Node, d float64, n int) {
	var readerWg sync.WaitGroup
	var writerWg sync.WaitGroup
	var minerWg sync.WaitGroup

	readerWg.Add(len(nodes))
	writerWg.Add(len(nodes))
	minerWg.Add(len(nodes))

	for i := 0; i < len(nodes); i++ {
		go threads.Reader(nodes[i], &readerWg)
		go threads.Writer(nodes[i], &writerWg)
		go threads.Miner(nodes[i], &minerWg, d, n)
	}
}

// e.g. go run main.go --nodes=5 --d=2.0 --n=4
func main() {
	fmt.Println("***** ******") // +1

	numOfNodes := flag.Uint("nodes", 6, "an unsigned int")
	miningDivisor := flag.Float64("d", 4, "a float")
	numOfSidelinks := flag.Int("n", 5, "an int")

	flag.Parse()
	fmt.Println("number of nodes: ", *numOfNodes)
	fmt.Println("mining divisor: ", *miningDivisor)
	fmt.Println("number of sidelinks: ", *numOfSidelinks)

	nodes := createNodes(*numOfNodes)
	initThreads(nodes, *miningDivisor, *numOfSidelinks)
	manager := api.Manager{Nodes: nodes}
	server := api.BlockchainServer{}

	go client.HandleInput()

	server.Run(&manager)
}
