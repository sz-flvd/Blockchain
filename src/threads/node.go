package threads

import (
	"sync"

	"krypto.blockchain/src/common"
)

/* Internal structure shared by miner subthreads */
type Internal struct {
	blockId   uint32
	blockPoW  []byte
	Timestamp int64
}

/* Structure shared by all miner subthreads */
type Node struct {
	NewRecordChannel chan []string
	readerChannel    chan Internal    // this is the channel on which Reader waits for information about newly mined blocks
	minerChannel     chan Internal    // Miner will inform Writer about a newly mined Block through this channel
	writerChannels   []*chan Internal // Writer will write to all of these channels when a new Block is mined by this Node
	/* 	Internal state of Node (naming may need to be adjusted;
	Reader will update this when a new block is mined outside of this Node
	and Miner will check if it still needs to mine the current Block by reading any updates in this struct) */
	state      Internal       // So i figure acces to this AND Chain has to be synced?
	Chain      []common.Block // Holds all Blocks mined in the current session
	lastBlock  *common.Block  // pointer to last Block mined in the current session (idk if this will be needed)
	chainMutex sync.Mutex
}

func Node_CreateNode(readerChannel chan Internal, writerChannels []*chan Internal) *Node {
	newNode := &Node{
		NewRecordChannel: make(chan []string),
		readerChannel:    readerChannel,
		minerChannel:     make(chan Internal),
		writerChannels:   writerChannels,
		state:            Internal{},
		Chain:            make([]common.Block, 0),
		chainMutex:       sync.Mutex{},
	}

	return newNode
}
