package drsm

import (
	"fmt"
    "log"
	"github.com/omec-project/MongoDBLibrary"
	"go.mongodb.org/mongo-driver/bson"
	//"math/rand"
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
	var cn int32 = 1
/*
	for {
		cn = rand.Int32(16000)
		_, found := d.globalChunkTbl[cn]
		if found == true {
			continue
		}
		log.Println("Found chunk Id block ", cn)
		break
	}
*/
	// Let's confirm if this gets updated in DB
	docId := fmt.Sprintf("chunkid-%d", cn)
	filter := bson.M{"_id": docId}
	update := bson.M{"_id": docId, "type": "chunk","podId": d.clientId.PodName}
	inserted := MongoDBLibrary.RestfulAPIPost(d.sharedPoolName, filter, update)
	if inserted != true {
		log.Printf("Adding chunk %v failed ",cn)
		err := fmt.Errorf("Ids not available")
		return nil, err
	}

	c := &Chunk{Id: cn}

	d.localChunkTbl[cn] = c

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
