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
	PodName string `bson:"podName,omitempty" json:"podName,omitempty"`
	PodIp   string `bson:"podIp,omitempty" json:"podIp,omitempty"`
}

type PodData struct {
	mu            sync.Mutex      `bsin:"-" json:"-"`
	PodId         PodId           `bson:"podId,omitempty" json:"podId,omitempty"`
	Timestamp     time.Time       `bson:"time,omitempty" json:"time,omitempty"`
	PrevTimestamp time.Time       `bson:"-" json:"-"`
	podChunks     map[int32]Chunk `bson:"-" json:"-"`
}

type PodHealth struct {
	Id      string  `bson:"_id,omitempty" json:"_id,omitempty"`
	podData PodData `bson:"podData,omitempty" json:"podData,omitempty"`
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
