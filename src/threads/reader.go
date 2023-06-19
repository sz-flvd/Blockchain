/*	General miner thread is divided into three parts:
		- Reader
		- Miner
		- Writer
	This is the part responsible for reading all updates from other threads
*/

package threads

import (
	"sync"

	"krypto.blockchain/src/common"
)

func Reader(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range node.readerChannelBlockMined {
		// read information about newly mined Blocks from other Nodes
		// this needs synchronization!!!
		node.state.blockId = i.blockId
		node.state.blockPoW = i.blockPoW
	}

	for {
		select {
		case newBlock := <-node.readerChannelBlockMined:
			continue
		case addedRecordData := <-node.readerChannelRecordAdd:
			senderId := addedRecordData.uint
			addedrecordPtr := addedRecordData.Record
			addedRecord := *(addedrecordPtr)
			// Send a confirmation to the sender
			*node.writerChannelsRecordConfirm[senderId] <- addedrecordPtr
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
				myNewRecord := common.Record{recordId, ts, content}
				awaiting := struct {
					common.Record
					uint
				}{
					Record: myNewRecord,
					uint:   1,
				}
				node.awaitingRecords = append(node.awaitingRecords, awaiting)
				for idx := uint(0); idx < node.networkSize; idx++ {
					if idx != node.index {
						*node.writerChannelsRecordAdd[idx] <- struct {
							*common.Record
							uint
						}{Record: &myNewRecord, uint: node.index}
					}
				}
			}
		case confirmedRecord := <-node.readerChannelRecordConfirm:
			confirmedRecordDerefed := *(confirmedRecord)
			content := confirmedRecordDerefed.Content
			_, foundIdx, doesContain := node.FindAwaitingRecord(content)
			if doesContain {
				// Increment confirmations for this record.
				node.awaitingRecords[foundIdx].uint++
				if node.awaitingRecords[foundIdx].uint >= node.networkSize { // All other nodes confirmed this record.
					// Pop this record from awaiting slice and push it into current block records.
					node.currentBlock.Records = append(node.currentBlock.Records, node.awaitingRecords[foundIdx].Record)
					node.awaitingRecords = append(node.awaitingRecords[:foundIdx], node.awaitingRecords[foundIdx+1:]...)
				}
			}
		}
	}
}
