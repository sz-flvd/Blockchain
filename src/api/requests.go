package api

import (
	"krypto.blockchain/src/common"
)

/* Request for adding a new record to blockchain */
type AddRecordRequest struct {
	Content   string    // arbitrary data
	Receivers []int  	// ids of nodes that are expected handle the request. Empty if no preference
}

// TODO docs
type AddRecordResponse struct {
	BlockId int    // id of block the record will be added to
	Index 	int    // index of record within block
}

type GetBlockResponse struct {
	BlockData common.Block  // Full information regarding block
}

type GetLocalBlockchainResponse struct {
	Blocks []common.Block  // the current blockchain of the specific node, including unconfirmed blocks
}