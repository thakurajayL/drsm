package drsm

import (
	"sync"
    "time"
)

type DbInfo struct {
	Url  string
	Name string
}

type PodId struct {
	PodName string `bson:"podName,omitempty`
	PodIp   string `bson:"podIp,omitempty`
}

type PodData struct {
	PodId      string                 `bson:"podId,omitempty"`
	Timestamp  time.Time              `bson:"time,omitempty"`
	PrevTimestamp  time.Time          `bson:"-"`
	podChunks  map[int32]Chunki       `bson:"-"`
	mu         sync.Mutex             `bsin:"-"`
}

type ChunkState int

const (
	Invalid ChunkState = iota + 1
	Owned
	PeerOwned
	Orphan
)

type Drsm struct {
	mu             sync.Mutex
	sharedPoolName string
	clientId       PodId
	db             DbInfo
	localChunkTbl  map[int32]*Chunk
	globalChunkTbl map[int32]*Chunk
	podMap         map[string]*PodData
	newPod         chan string
	podDown        chan string
	scanChunk      chan int32
}
