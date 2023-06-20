package api

import (
	"bytes"
	"errors"
	"strconv"

	"krypto.blockchain/src/threads"
)

type Manager struct {
	Nodes []*threads.Node
}

// Fun fact: I believe that return type (*something, error) is actually a monad 'Either'
func (manager *Manager) AddRecord(request AddRecordRequest) (*AddRecordResponse, error) {
	content := request.Content

	for _, i := range request.Receivers {
		if i < len(manager.Nodes) {
			manager.Nodes[i].NewRecordChannel <- content
		}
	}

	if len(request.Receivers) == 0 {
		for i := range manager.Nodes {
			manager.Nodes[i].NewRecordChannel <- content
		}
	}

	// TODO: 'The network should return the block number that will contain the record, even if
	// the block hasn’t been confirmed yet (see below). Each record should have a unique identifier
	// that should also be returned by the network. Block id may be part of the record id, i.e. it’s
	// enough to consider (block-index, index-within-block) as record id.'
	// Not sure, how to do this, honestly. For now, the function will return the pair:
	// (index of last mined block + 1, 0), but it will probably need to be changed
	return &AddRecordResponse{BlockId: len(manager.Nodes[0].Chain), Index: 0}, nil
}

func (manager *Manager) GetBlock(blockId int) (*GetBlockResponse, error) {
	if !isBlockConfirmed(manager, blockId) {
		return nil, errors.New("Block with id: " + strconv.Itoa(blockId) + " doesn't exist or is not confirmed yet")
	}

	return &GetBlockResponse{BlockData: manager.Nodes[0].Chain[blockId]}, nil
}

func (manager *Manager) GetLocalBlockchain(index int) (*GetLocalBlockchainResponse, error) {
	if index >= len(manager.Nodes) {
		return nil, errors.New("Incorrect node id: " + strconv.Itoa(index))
	}

	response := GetLocalBlockchainResponse{Blocks: manager.Nodes[index].Chain}

	return &response, nil
}

func (manager *Manager) GetBlockCount() int {
	count := 0
	for i := range manager.Nodes[0].Chain {

		if !isBlockConfirmed(manager, i) {
			return count
		}
		count++
	}

	return count
}

func isBlockConfirmed(manager *Manager, index int) bool {
	// Currently, we don't explicitely track, which blocks are confirmed
	// As such, this function assumes that the block is confirmed if it stored in all Nodes
	// in the same position and with the same PoW
	// Btw, I assume the Block ids are their index in the chain are the same thing, which might be pretty bold on my part
	if len(manager.Nodes[0].Chain) <= index {
		return false
	}
	block := manager.Nodes[0].Chain[index]

	for i := range manager.Nodes {
		if len(manager.Nodes[i].Chain) <= index {
			return false
		}

		if !bytes.Equal(manager.Nodes[i].Chain[index].PoW, block.PoW) {
			return false
		}
	}

	return true
}
