/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for actually mining block
*/

package threads

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"

	"krypto.blockchain/src/common"
)

var d float64
var n int

const (
	l = sha256.Size
)

func Miner(node *Node, wg *sync.WaitGroup, divisor float64, sidelinks int) {
	defer wg.Done()

	d = divisor
	n = sidelinks

	for record := range node.NewRecordChannel {
		// if !ok {
		// 	return
		// }
		records := []common.Record{record}

	scanAllRecords:
		for {
			select {
			case nT := <-node.NewRecordChannel:
				records = append(records, nT)
			case <-time.After(1 * time.Millisecond):
				break scanAllRecords
			}
		}

		data, newBlock := prepareBlockAndData(node, records)
		h := sha256.Sum256(data)
		mined := false

		b := make([]byte, l/2)
		var timestamp int64
		for {

			rand.Read(b)
			// pick random b
			token := sha256.Sum256(append(h[:], b...))
			timestamp = time.Now().UnixNano()
			fmt.Printf("Node %v calculating hash value = %v vs diff %v\n", node.index, TokenValue(token), math.Pow(2.0, -d))
			if TokenValue(token) < math.Pow(2.0, -d) {
				mined = true
				break
			} else if node.state.blockPoW != nil { // this needs some synchronization!!!
				// I guess it can work this way, that checking elements of the chain has to be synchronized
				// we may on the other hand consider creating separate thread that will provide synchronized RW actions on chain
				// through usage of select statement
				node.chainMutex.Lock()
				if node.state.blockId != node.lastBlock.Index {
				    node.chainMutex.Unlock()
					break
				}
				node.chainMutex.Unlock()
			}
			select {
			case newTransaction := <-node.NewRecordChannel:
				records = append(records, newTransaction)
			readerLoop:
				for {
					select {
					case nT := <-node.NewRecordChannel:
						records = append(records, nT)
					case <-time.After(1 * time.Millisecond):
						break readerLoop
					}
				}
				data, newBlock = prepareBlockAndData(node, records)
				h = sha256.Sum256(data)
				mined = false
			case <-time.After(1 * time.Millisecond):

			}
		}
		fmt.Printf("Node %v done \n", node.index)
		// We need something here to synchronize with all other nodes, that they accept out firsthood
		// i propose something like channel waiting for 8 messeges, if all are OK then accept
		// in case anyone does not accept we need to figure out some protocol
		// How about "The earliest timestamp -> lowest b -> lowest index"

		// Also we need to somehow check if anyone sent us this information here OR at any time if our chain is shorter than anyones else
		// And if we accept someones else block, we have to calculate hash with PoW to prove it

		if mined {
			newBlock.PoW = b
			newBlock.Timestamp = timestamp
			node.minerChannel <- Internal{
				blockId:   newBlock.Index,
				blockPoW:  b,
				Timestamp: timestamp,
			}
		} else {
			newBlock.PoW = node.state.blockPoW
			newBlock.Timestamp = node.state.Timestamp
		}

		node.Chain = append(node.Chain, newBlock)
		node.lastBlock = &node.Chain[len(node.Chain)-1]
	}
}

func TokenValue(token [l]byte) float64 {
	// This is BAD, but I don't know how to do this correctly.
	// We have 32 bytes of output and want to get 8 bit float...
	val := binary.BigEndian.Uint64(token[0:8])
	return float64(val) / float64(^uint(0))
}

func prepareBlockAndData(node *Node, records []common.Record) ([]byte, common.Block) {
	prevBlockHash := calcBlockHash(node.lastBlock)
	sideLinks := calcSidelinks(node.Chain)
	// records := createRecords(records)
	newBlock := common.Block{
		Index:       node.lastBlock.Index + 1,
		MainHash:    prevBlockHash[:],
		ExtraHashes: sideLinks,
		Records:     records,
	}
	data := make([]byte, 0)
	data = append(data, byte(newBlock.Index))
	data = append(data, newBlock.MainHash...)
	for _, hash := range newBlock.ExtraHashes {
		data = append(data, hash...)
	}

	return data, newBlock
}

func calcBlockHash(b *common.Block) []byte {
	data := make([]byte, 0)
	data = append(data, byte(b.Index))
	data = append(data, b.MainHash...)
	for _, hash := range b.ExtraHashes {
		data = append(data, hash...)
	}

	// timestampHolder := make([]byte, 0)
	// binary.LittleEndian.PutUint64(timestampHolder, uint64(b.Timestamp))
	// data = append(data, timestampHolder...)

	for _, record := range b.Records {
		data = append(data, record.ToBytes()...)
	}

	h := sha256.Sum256(data)
	res := sha256.Sum256(append(h[:], b.PoW...))
	return res[:]
}

func calcSidelinks(chain []common.Block) [][]byte {
	sidelinks := make([][]byte, 0)
	if len(chain) <= n {
		for _, block := range chain {
			sidelinks = append(sidelinks, calcBlockHash(&block)[:])
		}

		return sidelinks
	}

	I := len(chain)

	prevHash := calcBlockHash(&chain[len(chain)-1])
	sidelinks = append(sidelinks, prevHash)

	indexes := make([]int, 0)

	x := prevHash
	howManyBytes := int(math.Ceil(float64(n) / float64(8)))
	xIntVal := 0

	for i := len(x) - howManyBytes; i < len(x); i++ {
		xIntVal *= 258
		xIntVal += int(x[i])
	}

	indexes = append(indexes, xIntVal%(I-1))
	for j := 1; j < n; j++ {
		byteHelper := make([]byte, 0)
		binary.LittleEndian.PutUint64(byteHelper, uint64(j))
		xj := sha256.Sum256(append(x, byteHelper...))

		xjIntVal := 0

		for i := len(prevHash) - howManyBytes; i < len(prevHash); i++ {
			xjIntVal *= 258
			xjIntVal += int(xj[i])
		}

		nj := xjIntVal % (I - (j + 1))

		for k, nk := range indexes {
			if nj == nk {
				nj = I - j + k - 1
				break
			}
		}

		indexes = append(indexes, nj)
	}

	for _, index := range indexes {
		sidelinks = append(sidelinks, calcBlockHash(&chain[index]))
	}

	return sidelinks
}
