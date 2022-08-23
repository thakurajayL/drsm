package drsm

type DbInfo struct {
	Url		string
	Name	string
}

type PodId struct {
	PodName  string
	PodIp	 string
}

type ChunkState int

const (
	Invalid ChunkState = iota + 1
	Owned
	PeerOwned
	Orphan
)

type PodHealth struct {
	mu      		sync.Mutex
	prevCount 		Int32
	currentCount 	Int32
}

type Drsm struct {
	mu      		sync.Mutex
	sharedPoolName	string
	clientId		PodId
	db				DbInfo
	localChunkTbl	map[Int32]*Chunk
	globalChunkTbl	map[Int32]*Chunk
	podMap          map[string]*PodHealth
}


