package drsm

import (
	"fmt"
	"github.com/omec-project/MongoDBLibrary"
)

func InitDRSM(sharedPoolName string, myid PodId, db DbInfo) (*Drsm, error) {
	d := &Drsm{sharedPoolName: sharedPoolName,
		clientId: myid,
		db:       db}
	d.localChunkTbl = make(map[int]*Chunk)
	d.globalChunkTbl = make(map[int]*Chunk)
	d.newPod = make(chan string, 10)

	//connect to DB
	MongoDBLibrary.SetMongoDB(db.Name, db.Url)
	handleDbUpdates(d)
	go startDiscovery(d)
	return d, nil
}

func (d *Drsm) AllocateIntID(sharedPoolName string) (int, error) {
	for k, c := range d.localChunkTbl {
		if len(c.FreeIds) > 0 {
			return c.AllocateIntID()
		}
	}
	c, err := GetNewChunk(d)
	if err {
		err := fmt.Errorf("Ids not available")
		return 0, err
	}
	return c.AllocateIntID(), nil
}

func (d *Drsm) ReleaseIntID(sharedPoolName string, id int) (error) {
	chunkId := id >> 10
	chunk, found := d.localChunkTbl[chunkId]
	if found == true {
		c.ReleaseIntID(id)
		return nil
	}
	err := fmt.Errorf("Unknown Id")
	return err
}

func (d *Drsm) FindOwnerIntID(sharedPoolName string, id int) (string, error) {
	chunkId := id >> 10
	i := id & 0x3ff
	chunk, found := d.localChunkTbl[chunkId]
	if found == true {
		return chunk.GetOwner()
	}
	err := fmt.Errorf("Unknown Id")
	return "", err
}
