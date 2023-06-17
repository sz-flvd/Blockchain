/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for actually mining block
*/

package threads

import (
	"crypto/sha256"
	"math"
	"sync"

	"krypto.blockchain/src/common"
)

const (
	l = sha256.Size
	d = 1.0 // this needs to be adjusted and set as a runtime parametre
)

func Miner(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		b, ok := <- node.NewRecordChannel
		if (!ok) {
			return
		}

		data := node.lastBlock.MainHash
		for _, extra := range node.lastBlock.ExtraHashes {
			data = append(data, extra...)
		}
		h := sha256.Sum256(data)
		mined := false

		for {
			// pick random b

			token := sha256.Sum256(append(h[:], b...))

			if tokenValue(token) < math.Pow(2.0, float64(l))/d {
				mined = true
				break
			} else if node.state.blockId == node.lastBlock.Index && node.state.blockPoW != nil { // this needs some synchronization!!!
				break
			}
		}

		if mined {
			node.lastBlock.PoW = b
			node.minerChannel <- Internal{
				blockId:  node.lastBlock.Index,
				blockPoW: b}
		} else {
			node.lastBlock.PoW = node.state.blockPoW
		}

		mined = false

		node.Blocks = append(node.Blocks, common.Block{
			// properly add new block to the Blockchain
		})
		node.lastBlock = &node.Blocks[len(node.Blocks)-1]
	}
}

func tokenValue(token [l]byte) float64 {
	val := 0.0

	for i, b := range token {
		val += float64(b) * math.Pow(2.0, float64(i))
	}

	return val
}
