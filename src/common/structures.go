package common

import "encoding/binary"

/* Record structure */
type Record struct {
	Index     uint32 // Record number within its corresponding Block
	Timestamp int64  // timestamp of adding the record into the Block
	Content   string // arbitrary data of Record
}

/* Block structure */
type Block struct {
	Index       uint32   // index of the Block
	Timestamp   int64    // timestamp of computing PoW
	MainHash    []byte   // hash of the previous Block
	ExtraHashes [][]byte // hashes of n arbitrarily chosen Blocks
	PoW         []byte   // proof of work of the Block
	Records     []Record // list of Records within the Block
}

func (r *Record) UpdateEarlierTimestamp(ts int64) {
	if ts < r.Timestamp {
		r.Timestamp = ts
	}
}

func (r *Record) ToBytes() []byte {
	data := make([]byte, 0)

	byteHelper := make([]byte, 0)
	binary.LittleEndian.PutUint64(byteHelper, uint64(r.Index))
	data = append(data, byteHelper...)

	byteHelper = make([]byte, 0)
	binary.LittleEndian.PutUint64(byteHelper, uint64(r.Timestamp))
	data = append(data, byteHelper...)

	data = append(data, []byte(r.Content)...)

	return data
}
