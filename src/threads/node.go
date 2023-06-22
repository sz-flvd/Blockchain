package threads

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	// "crypto/sha256"
	// "crypto/x509"
	// "encoding/asn1"
	"fmt"
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

type MessageBool struct {
	b  bool
	id uint
}

type MessageBlock struct {
	block     common.Block
	id        uint
	publicKey ecdsa.PublicKey
	signature []byte
}

/* Structure shared by all miner subthreads */
type Node struct {
	index                           uint
	networkSize                     uint
	NewRecordChannel                chan common.Record
	readerChannelBlockMined         chan MessageBlock // this is the channel on which Reader waits for information about newly mined blocks
	readerChannelBlockConfirmation  chan MessageBool
	readerChannelRecordAdd          chan RecordAdd      // Reader gets info about newly added records.
	readerChannelRecordConfirm      chan *common.Record // Reader gets confirmation from other Nodes about his newly added record.
	minerChannel                    chan MessageBlock   // Miner will inform Writer about a newly mined Block through this channel
	rejectChannel                   chan MessageBlock
	acceptChannel                   chan MessageBlock
	writerChannelsBlockMined        []chan MessageBlock // Writer will write to all of these channels when a new Block is mined by this Node
	writerChannelsBlockConfirmation []chan MessageBool
	writerChannelsRecordAdd         []chan RecordAdd
	writerChannelsRecordConfirm     []chan *common.Record
	/* 	Internal state of Node (naming may need to be adjusted;
	Reader will update this when a new block is mined outside of this Node
	and Miner will check if it still needs to mine the current Block by reading any updates in this struct) */
	state             Internal       // So i figure access to this AND Chain has to be synced?
	Chain             []common.Block // Holds all Blocks mined in the current session
	lastBlock         *common.Block  // pointer to last Block mined in the current session (idk if this will be needed)
	currentBlock      common.Block   // block we're currently calculating PoW on.
	currentAcceptance uint
	chainMutex        sync.Mutex
	recordMutex       sync.Mutex
	// hasNewConfirmedRecords bool
	awaitingRecords []struct {
		common.Record
		uint
	}
	minerStop          bool
	waitingForApproval bool
	// Currency wallet
	privateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func Node_CreateNode(
	index uint,
	networkSize uint,
	readerChannelBlockMined chan MessageBlock,
	readerChannelBlockConfirmation chan MessageBool,
	readerChannelRecordAdd chan RecordAdd,
	readerChannelRecordConfirm chan *common.Record,
	writerChannelsBlockMined []chan MessageBlock,
	writerChannelsBlockConfirmation []chan MessageBool,
	writerChannelsRecordAdd []chan RecordAdd,
	writerChannelsRecordConfirm []chan *common.Record,
) *Node {
	// Cryptocurrency
	walletKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating ECDSA key pair:", err)
		return nil
	}
	// fmt.Println(string((*walletKey.PublicKey.X).String()))

	newNode := &Node{
		index:                           index,
		networkSize:                     networkSize,
		NewRecordChannel:                make(chan common.Record, 8),
		readerChannelBlockMined:         readerChannelBlockMined,
		readerChannelBlockConfirmation:  readerChannelBlockConfirmation,
		readerChannelRecordAdd:          readerChannelRecordAdd,
		readerChannelRecordConfirm:      readerChannelRecordConfirm,
		minerChannel:                    make(chan MessageBlock),
		rejectChannel:                   make(chan MessageBlock),
		acceptChannel:                   make(chan MessageBlock),
		writerChannelsBlockConfirmation: writerChannelsBlockConfirmation,
		writerChannelsBlockMined:        writerChannelsBlockMined,
		writerChannelsRecordAdd:         writerChannelsRecordAdd,
		writerChannelsRecordConfirm:     writerChannelsRecordConfirm,
		state:                           Internal{},
		Chain:                           make([]common.Block, 0),
		chainMutex:                      sync.Mutex{},
		recordMutex:                     sync.Mutex{},

		// hasNewConfirmedRecords:      false,
		awaitingRecords: make([]struct {
			common.Record
			uint
		}, 0),
		privateKey: walletKey,
		PublicKey:  &walletKey.PublicKey,
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

// Cryptocurrency
func (node *Node) AccessPrivateKey() *ecdsa.PrivateKey {
	return node.privateKey
}
