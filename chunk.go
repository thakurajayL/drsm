package drsm

import (
	"fmt"
	"github.com/omec-project/MongoDBLibrary"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func (c *chunk) GetOwner() *PodId {
	return &c.Owner

}

func (d *Drsm) GetNewChunk() (*chunk, error) {
	// Get new Chunk
	// We got to allocate new Chunk. We should select
	// probable chunk number

	log.Println("Allocate new chunk ")
	// 14 bits --- 1,2,4,8,16
	var cn int32 = 1
	for {
		for {
			cn = rand.Int31n(d.chunkIdRange)
			_, found := d.globalChunkTbl[cn]
			if found == true {
				continue
			}
			log.Println("Found chunk Id block ", cn)
			break
		}
		// Let's confirm if this gets updated in DB
		docId := fmt.Sprintf("chunkid-%d", cn)
		filter := bson.M{"_id": docId}
		update := bson.M{"_id": docId, "type": "chunk", "podId": d.clientId.PodName}
		inserted := MongoDBLibrary.RestfulAPIPostOnly(d.sharedPoolName, filter, update)
		if inserted != true {
			log.Printf("Adding chunk %v failed. Retry again ", cn)
			continue
		}
		break
	}

	log.Printf("Adding chunk %v success ", cn)
	c := &chunk{Id: cn}
	c.AllocIds = make(map[int32]bool)
	var i int32
	for i = 0; i < 1000; i++ {
		c.FreeIds = append(c.FreeIds, i)
	}

	c.resourceValidCb = d.resourceValidCb
	d.localChunkTbl[cn] = c

	// add Ids to freeIds
	return c, nil
}

func (c *chunk) AllocateIntID() int32 {
	if len(c.FreeIds) == 0 {
		log.Println("FreeIds in chunk 0")
		return 0
	}
	id := c.FreeIds[len(c.FreeIds)-1]
	c.FreeIds = c.FreeIds[:len(c.FreeIds)-1]
	return (c.Id << 10) | id
}

func (c *chunk) ReleaseIntID(id int32) {
	var i int32
	i = id & 0x3ff
	c.FreeIds = append(c.FreeIds, i)
}

func getChunIdFromDocId(id string) int32 {
	z := strings.Split(id, "-")
	if len(z) == 2 && z[0] == "chunkid" {
		cid, _ := strconv.ParseInt(z[1], 10, 32)
		c := int32(cid)
		return c
	}
	return 0
}
func isChunkDoc(id string) bool {
	z := strings.Split(id, "-")
	if len(z) == 2 && z[0] == "chunkid" {
		return true
	}
	return false
}
