package common

import (
	"crypto/sha256"
	"math"
	"time"
)

const (
	l = sha256.Size
	d = 1.0 // this needs to be adjusted and set as a runtime parametre
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

func (block *Block) ProofOfWork() {
	b := make([]byte, 0)
	data := block.MainHash
	for _, extra := range block.ExtraHashes {
		data = append(data, extra...)
	}
	h := sha256.Sum256(data)

	for {
		// pick random b
		token := sha256.Sum256(append(h[:], b...))

		if tokenValue(token) < math.Pow(2.0, float64(l))/d {
			break
		}
	}

	block.PoW = b
}

func tokenValue(token [l]byte) float64 {
	val := 0.0

	for i, b := range token {
		val += float64(b) * math.Pow(2.0, float64(i))
	}

	return val
}
