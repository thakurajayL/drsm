package drsm

import (
	"github.com/omec-project/MongoDBLibrary"
	"log"
	"math/rand"
	"sync"
	"time"
)

type DbInfo struct {
	Url  string
	Name string
}

type ChunkState int

const (
	Invalid ChunkState = iota + 1
	Owned
	PeerOwned
	Orphan
)

type Chunk struct {
	Id       int32
	Owner    PodId
	State    ChunkState
	FreeIds  []int32
	AllocIds map[int32]bool
}

type PodId struct {
	PodName string `bson:"podName,omitempty" json:"podName,omitempty"`
	PodIp   string `bson:"podIp,omitempty" json:"podIp,omitempty"`
}

type PodData struct {
	mu            sync.Mutex       `bsin:"-" json:"-"`
	PodId         PodId            `bson:"podId,omitempty" json:"podId,omitempty"`
	Timestamp     time.Time        `bson:"time,omitempty" json:"time,omitempty"`
	PrevTimestamp time.Time        `bson:"-" json:"-"`
	podChunks     map[int32]*Chunk `bson:"-" json:"-"` // chunkId to Chunk
}

type Drsm struct {
	mu             sync.Mutex
	sharedPoolName string
	clientId       PodId
	db             DbInfo
	localChunkTbl  map[int32]*Chunk    // chunkid to chunk
	globalChunkTbl map[int32]*Chunk    // chunkid to chunk
	podMap         map[string]*PodData // podId to podData
	newPod         chan string
	podDown        chan string
	scanChunk      chan int32
}

func (d *Drsm) ConstuctDrsm() {
	d.localChunkTbl = make(map[int32]*Chunk)
	d.globalChunkTbl = make(map[int32]*Chunk)
	d.podMap = make(map[string]*PodData)
	d.newPod = make(chan string, 10)
	d.podDown = make(chan string, 10)
	d.scanChunk = make(chan int32, 10)
	t := time.Now().UnixNano()
	rand.Seed(t)
	//connect to DB
	MongoDBLibrary.SetMongoDB(d.db.Name, d.db.Url)
	log.Println("SetMongoDB done ", d.db.Name)
	go d.handleDbUpdates()
	go d.startDiscovery()
	go d.podDownDetected()
}
