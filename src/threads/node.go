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
	index                   uint
	networkSize             uint
	NewRecordChannel        chan []string
	readerChannelBlockMined chan Internal // this is the channel on which Reader waits for information about newly mined blocks
	readerChannelRecordAdd  chan struct {
		*common.Record
		uint
	} // Reader gets info about newly added records.
	readerChannelRecordConfirm chan *common.Record // Reader gets confirmation from other Nodes about his newly added record.
	minerChannel               chan Internal       // Miner will inform Writer about a newly mined Block through this channel
	writerChannelsBlockMined   []*chan Internal    // Writer will write to all of these channels when a new Block is mined by this Node
	writerChannelsRecordAdd    []*chan struct {
		*common.Record
		uint
	}
	writerChannelsRecordConfirm []*chan *common.Record
	/* 	Internal state of Node (naming may need to be adjusted;
	Reader will update this when a new block is mined outside of this Node
	and Miner will check if it still needs to mine the current Block by reading any updates in this struct) */
	state           Internal       // So i figure access to this AND Chain has to be synced?
	Chain           []common.Block // Holds all Blocks mined in the current session
	lastBlock       *common.Block  // pointer to last Block mined in the current session (idk if this will be needed)
	currentBlock    common.Block   // block we're currently calculating PoW on.
	chainMutex      sync.Mutex
	awaitingRecords []struct {
		common.Record
		uint
	}
}

func Node_CreateNode(
	index uint,
	networkSize uint,
	readerChannelBlockMined chan Internal,
	readerChannelRecordAdd chan struct {
		*common.Record
		uint
	},
	readerChannelRecordConfirm chan *common.Record,
	writerChannelsBlockMined []*chan Internal,
	writerChannelsRecordAdd []*chan struct {
		*common.Record
		uint
	},
	writerChannelsRecordConfirm []*chan *common.Record,
) *Node {
	newNode := &Node{
		index:                       index,
		networkSize:                 networkSize,
		NewRecordChannel:            make(chan []string),
		readerChannelBlockMined:     readerChannelBlockMined,
		readerChannelRecordAdd:      readerChannelRecordAdd,
		readerChannelRecordConfirm:  readerChannelRecordConfirm,
		minerChannel:                make(chan Internal),
		writerChannelsBlockMined:    writerChannelsBlockMined,
		writerChannelsRecordAdd:     writerChannelsRecordAdd,
		writerChannelsRecordConfirm: writerChannelsRecordConfirm,
		state:                       Internal{},
		Chain:                       make([]common.Block, 0),
		chainMutex:                  sync.Mutex{},
		awaitingRecords: make([]struct {
			common.Record
			uint
		}, 0),
	}

	return newNode
}

func (node *Node) IndexOfRecordContainingContent(content string) (*common.Record, uint, bool) {
	for idx, myR := range node.currentBlock.Records {
		if myR.Content == content {
			return &myR, uint(idx), true
		}
	}

	for idx, myR := range node.awaitingRecords {
		if myR.Content == content {
			return &myR.Record, uint(idx), true
		}
	}

	return nil, 0, false
}
