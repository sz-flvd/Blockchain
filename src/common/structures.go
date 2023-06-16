package common

import (
	"time"
)

/* Record structure */
type Record struct {
	Index     uint32    // Record number within its corresponding Block
	Timestamp time.Time // timestamp of adding the record into the Block
	Content   []byte    // arbitrary data of Record
}

/* Block structure */
type Block struct {
	Index       uint32    // index of the Block
	Timestamp   time.Time // timestamp of computing PoW
	MainHash    []byte    // hash of the previous Block
	ExtraHashes [][]byte  // hashes of n arbitrarily chosen Blocks
	PoW         []byte    // proof of work of the Block
	Records     []Record  // list of Records within the Block
}
