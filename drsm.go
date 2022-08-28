package drsm

import (
	"github.com/omec-project/MongoDBLibrary"
	"log"
	"math/rand"
	"sync"
	"time"
)

type chunkState int

const (
	Invalid chunkState = iota + 1
	Owned
	PeerOwned
	Orphan
	Scanning
)

type chunk struct {
	Id              int32
	Owner           PodId
	State           chunkState
	FreeIds         []int32
	AllocIds        map[int32]bool
	ScannedIds      []int32
	stopScan        chan bool
	resourceValidCb func(int32) bool
}

type podData struct {
	mu            sync.Mutex       `bsin:"-" json:"-"`
	PodId         PodId            `bson:"podId,omitempty" json:"podId,omitempty"`
	Timestamp     time.Time        `bson:"time,omitempty" json:"time,omitempty"`
	PrevTimestamp time.Time        `bson:"-" json:"-"`
	podChunks     map[int32]*chunk `bson:"-" json:"-"` // chunkId to Chunk
}

type Drsm struct {
	mu              sync.Mutex
	sharedPoolName  string
	clientId        PodId
	db              DbInfo
	mode            DrsmMode
	resIdSize       int32
	localChunkTbl   map[int32]*chunk    // chunkid to chunk
	globalChunkTbl  map[int32]*chunk    // chunkid to chunk
	podMap          map[string]*podData // podId to podData
	podDown         chan string
	scanChunks      map[int32]*chunk
	chunkIdRange    int32
	resourceValidCb func(int32) bool
}

func (d *Drsm) ConstuctDrsm(opt *Options) {
	if opt != nil {
		d.mode = opt.Mode
		log.Println("drsm mode set to ", d.mode)
		if opt.ResIdSize > 0 {
			d.resIdSize = opt.ResIdSize
		} else {
			d.resIdSize = 24
		}
		d.resourceValidCb = opt.ResourceValidCb
	}
	d.chunkIdRange = 1 << (d.resIdSize - 10)
	log.Println("ChunkId in the range of 0 to %v ", d.chunkIdRange)
	d.localChunkTbl = make(map[int32]*chunk)
	d.globalChunkTbl = make(map[int32]*chunk)
	d.podMap = make(map[string]*podData)
	d.podDown = make(chan string, 10)
	d.scanChunks = make(map[int32]*chunk)
	t := time.Now().UnixNano()
	rand.Seed(t)

	//connect to DB
	MongoDBLibrary.SetMongoDB(d.db.Name, d.db.Url)
	log.Println("SetMongoDB done ", d.db.Name)

	go d.handleDbUpdates()
	go d.punchLiveness()
	go d.podDownDetected()
	go d.checkAllChunks()
}
