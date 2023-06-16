/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for writing updates to other threads
*/

package threads

import "sync"

func Writer(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range node.minerChannel {
		// broadcast information about mined Block to all Reader threads
		for _, c := range node.writerChannels {
			(*c) <- i
		}
	}
}
