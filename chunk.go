package drsm

import (
	"fmt"
	"github.com/omec-project/MongoDBLibrary"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
)

type Chunk struct {
	Id       int32
	Owner    PodId
	State    ChunkState
	FreeIds  []int32
	AllocIds map[int32]bool
}

func (c *Chunk) GetOwner() PodId {
	return c.Owner

}

func GetNewChunk(d *Drsm) (*Chunk, error) {
	// Get new Chunk
	// We got to allocate new Chunk. We should select
	// probable chunk number

	log.Println("Allocate new chunk ")
	// 14 bits --- 1,2,4,8,16
	var cb int32
	for {
		cn = rand.int32(16000)
		_, found := d.globalChunkTbl[cn]
		if found == true {
			continue
		}
		log.Println("Found chunk Id block ", cn)
		break
	}
	// Let's confirm if this gets updated in DB
	docId := fmt.Sprintf("chunkid-%s", cn)
	filter := bson.M{"_id": docId}
	update := bson.D{{"_id", docId}, {"type", "chunk"}, {"podId", d.clientId.PodName}}
	_, err := MongoDBLibrary.RestfulAPIPost(d.sharedPoolName, filter, update)
	if err != nil {
		log.Println("put data failed : ", err)
		return
	}

	c := &Chunk{Id: cn}

	d.localChunkTbl[id] = c

	// add Ids to freeIds
	return c, nil
}

func (c *Chunk) AllocateIntID() int32 {
	id := c.FreeIds[len(c.FreeIds)-1]
	c.FreeIds = c.FreeIds[:len(c.FreeIds)-1]
	return (c.Id << 10) | id
}

func (c *Chunk) ReleaseIntID(id int32) {
	var i int32
	i = id & 0x3ff
	c.FreeIds = append(c.FreeIds, i)
}
