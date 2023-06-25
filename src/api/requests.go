package api

import (
	// "crypto/ecdsa"

	"krypto.blockchain/src/common"
	"krypto.blockchain/src/threads"
)

/* Request for adding a new record to blockchain */
type AddRecordRequest struct {
	Content   []string // arbitrary data
	Receivers []int    // ids of nodes that are expected handle the request. Empty if no preference
}

type AddRecordResponse struct {
	BlockId int // id of block the record will be added to
	Index   int // index of record within block
}

type GetBlockResponse struct {
	BlockData common.Block // Full information regarding block
}

type GetLocalBlockchainResponse struct {
	Blocks []common.Block // the current blockchain of the specific node, including unconfirmed blocks
}

// Cryptocurrency
// type GetPublicKeyResponse struct {
// 	PublicKey ecdsa.PublicKey
// }

// type GetPrivateKeyResponse struct {
// 	PublicKey ecdsa.PrivateKey
// }

type GetNodeResponse struct { // XD
	NodeResponce *threads.Node
}
