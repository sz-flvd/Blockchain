/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for reading all updates from other threads
*/

package threads

import "sync"

func Reader(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range node.readerChannel {
		// read information about newly mined Blocks from other Nodes
		// this needs synchronization!!!
		node.state.blockId = i.blockId
		node.state.blockPoW = i.blockPoW
	}
}
