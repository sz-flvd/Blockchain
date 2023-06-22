/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for writing updates to other threads
*/

package threads

import (
	"sync"
)

func Writer(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	// for i := range node.minerChannel {
	// 	// broadcast information about mined Block to all Reader threads
	// 	for _, c := range node.writerChannelsBlockMined {
	// 		(*c) <- i
	// 	}
	// }
	// fmt.Println("Writer working?")
	for {
		select {
		case message := <-node.rejectChannel:
			node.writerChannelsBlockConfirmation[message.id] <- MessageBool{b: false, id: node.index}
		case message := <-node.acceptChannel:
			// fmt.Printf("Node %v Sendinnnnng conf to node %v\n", node.index, message.id)

			node.writerChannelsBlockConfirmation[message.id] <- MessageBool{b: true, id: node.index}
		case block := <-node.minerChannel:
			node.chainMutex.Lock()
			node.waitingForApproval = true
			for i, nodeChannel := range node.writerChannelsBlockMined {
				if i != int(node.index) {
					// fmt.Printf("Node %v SENDING %v\n", node.index, block.PoW)

					nodeChannel <- MessageBlock{block: block, id: node.index}
				}
			}
			node.chainMutex.Unlock()
		}

	}
}
