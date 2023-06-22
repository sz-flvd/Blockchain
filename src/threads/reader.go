/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for reading all updates from other threads
*/

package threads

import (
	"crypto/sha256"
	"math"
	"sync"

	"krypto.blockchain/src/common"
)

func Reader(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	// for i := range node.readerChannelBlockMined {
	// 	// read information about newly mined Blocks from other Nodes
	// 	// this needs synchronization!!!

	// }

	for {
		select {
		case addedRecordData := <-node.readerChannelRecordAdd:
			senderId := addedRecordData.sender
			addedrecordPtr := addedRecordData.record
			addedRecord := *(addedrecordPtr)
			// Send a confirmation to the sender
			node.writerChannelsRecordConfirm[senderId] <- addedrecordPtr
			// Check content
			recordId := addedRecord.Index
			ts := addedRecord.Timestamp
			content := addedRecord.Content
			foundPtr, _, doesContain := node.FindRecordContainingContent(content)
			if doesContain {
				// Update the timestamp to the earliest
				foundPtr.UpdateEarlierTimestamp(ts)
			} else {
				// Add a new Record with the same data to my structure and post for confirmations.
				myNewRecord := common.Record{
					Index:     recordId,
					Timestamp: ts,
					Content:   content,
				}
				awaiting := struct {
					common.Record
					uint
				}{
					Record: myNewRecord,
					uint:   1,
				}
				node.chainMutex.Lock()
				node.awaitingRecords = append(node.awaitingRecords, awaiting)
				node.chainMutex.Unlock()
				for idx := uint(0); idx < node.networkSize; idx++ {
					if idx != node.index {
						node.writerChannelsRecordAdd[idx] <- RecordAdd{record: &myNewRecord, sender: node.index}
					}
				}
			}
		case confirmedRecord := <-node.readerChannelRecordConfirm:
			confirmedRecordDerefed := *(confirmedRecord)
			content := confirmedRecordDerefed.Content
			node.recordMutex.Lock()
			_, foundIdx, doesContain := node.FindAwaitingRecord(content)
			if doesContain {
				// Increment confirmations for this record.
				node.awaitingRecords[foundIdx].uint++
				if node.awaitingRecords[foundIdx].uint >= node.networkSize { // All other nodes confirmed this record.
					// Pop this record from awaiting slice and push it into current block records.
					node.currentBlock.Records = append(node.currentBlock.Records, node.awaitingRecords[foundIdx].Record)
					node.awaitingRecords = append(node.awaitingRecords[:foundIdx], node.awaitingRecords[foundIdx+1:]...)
					// node.hasNewConfirmedRecords = true
					node.NewRecordChannel <- node.currentBlock.Records[len(node.currentBlock.Records)-1]
				}
			}
			node.recordMutex.Unlock()
		case message := <-node.readerChannelBlockMined:
			newBlock := message.block
			node.chainMutex.Lock()
			if node.waitingForApproval {
				node.acceptChannel <- message

			}
			if verifyBlock(&newBlock) {

				if compareHash(newBlock.MainHash, node.currentBlock.MainHash) {

					// fmt.Printf("Node %v Im here\n", node.index)
					node.Chain = append(node.Chain, newBlock)
					node.lastBlock = &node.Chain[len(node.Chain)-1]
					node.minerStop = true

					notAddedRecords := make([]common.Record, 0)

				outerLoop:
					for _, rec := range node.currentBlock.Records {

						for _, addRec := range newBlock.Records {
							if rec.Index == addRec.Index {
								break outerLoop
							}
						}
						notAddedRecords = append(notAddedRecords, rec)
					}

					node.currentBlock.Records = notAddedRecords
					for _, rec := range notAddedRecords {
						node.NewRecordChannel <- rec

					}
				}
				// fmt.Printf("Node %v Sent confirmation\n", node.index)

				node.acceptChannel <- message
			} else {
				node.rejectChannel <- message
			}
			node.chainMutex.Unlock()

		case accepted := <-node.readerChannelBlockConfirmation:
			node.chainMutex.Lock()
			// fmt.Printf("Node %v got something\n", node.index)

			if !accepted.b {
				continue
			}

			node.currentAcceptance++
			// fmt.Printf("Node %v got acceptance %v\n", node.index, node.currentAcceptance)

			if node.currentAcceptance == node.networkSize-1 {
				// fmt.Printf("Node %v all accepted\n", node.index)

				node.Chain = append(node.Chain, node.currentBlock)
				node.lastBlock = &node.Chain[len(node.Chain)-1]
				node.currentAcceptance = 0
			}
			node.chainMutex.Unlock()
		}

	}
}

func verifyBlock(block *common.Block) bool {
	data := make([]byte, 0)
	data = append(data, byte(block.Index))
	data = append(data, block.MainHash...)
	for _, hash := range block.ExtraHashes {
		data = append(data, hash...)
	}

	h := sha256.Sum256(data)
	token := sha256.Sum256(append(h[:], block.PoW...))

	return TokenValue(token) < math.Pow(2.0, -d)

}

func compareHash(h1, h2 []byte) bool {
	for i := range h1 {
		if h1[i] != h2[i] {
			return false
		}
	}

	return true
}
