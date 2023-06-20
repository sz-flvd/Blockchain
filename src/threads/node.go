package threads

import (
	"sync"
    "time"
    
	"krypto.blockchain/src/common"
)

/* Internal structure shared by miner subthreads */
type Internal struct {
	blockId   uint32
	blockPoW  []byte
	Timestamp int64
}

type RecordAdd struct {
    record *common.Record
    sender uint
}

/* Structure shared by all miner subthreads */
type Node struct {
	index                   uint
	networkSize             uint
	NewRecordChannel        chan []string
	readerChannelBlockMined chan Internal // this is the channel on which Reader waits for information about newly mined blocks
	readerChannelRecordAdd  chan RecordAdd // Reader gets info about newly added records.
	readerChannelRecordConfirm chan *common.Record // Reader gets confirmation from other Nodes about his newly added record.
	minerChannel               chan Internal       // Miner will inform Writer about a newly mined Block through this channel
	writerChannelsBlockMined   []*chan Internal    // Writer will write to all of these channels when a new Block is mined by this Node
	writerChannelsRecordAdd    []*chan RecordAdd
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
	readerChannelRecordAdd chan RecordAdd,
	readerChannelRecordConfirm chan *common.Record,
	writerChannelsBlockMined []*chan Internal,
	writerChannelsRecordAdd []*chan RecordAdd,
	writerChannelsRecordConfirm []*chan *common.Record,
) *Node {
	newNode := &Node{
		index:                       index,
		networkSize:                 networkSize,
		NewRecordChannel:            make(chan []string, 8),
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
	
	genesisBlock := common.Block{
    		Index:       0,
    		Timestamp:   time.Now().UnixNano(),
    		MainHash:    make([]byte, 32),
    		ExtraHashes: make([][]byte, 0),
    		PoW:         make([]byte, 0),
    		Records:     make([]common.Record, 0),
    	}
    
    	newNode.Chain = append(newNode.Chain, genesisBlock)
    	newNode.lastBlock = &genesisBlock

	return newNode
}

func (node *Node) FindAwaitingRecord(content string) (*struct {
	common.Record
	uint
}, uint, bool) {
	for idx, myR := range node.awaitingRecords {
		r := myR.Record
		if r.Content == content {
			return &myR, uint(idx), true
		}
	}

	return nil, 0, false
}

func (node *Node) FindRecordContainingContent(content string) (*common.Record, uint, bool) {
	for idx, myR := range node.currentBlock.Records {
		if myR.Content == content {
			return &myR, uint(idx), true
		}
	}

	found, idx, doesContain := node.FindAwaitingRecord(content)
	if doesContain {
		foundRecord := (*found).Record
		return &foundRecord, idx, true
	}

	return nil, 0, false
}
