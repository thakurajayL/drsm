package drsm

import (
	"fmt"
	"log"
)

func InitDRSM(sharedPoolName string, myid PodId, db DbInfo) (*Drsm, error) {
	log.Println("****MY ID ", myid)

	d := &Drsm{sharedPoolName: sharedPoolName,
		clientId: myid,
		db:       db}

	d.ConstuctDrsm()

	return d, nil
}

func (d *Drsm) AllocateIntID(sharedPoolName string) (int32, error) {
	for _, c := range d.localChunkTbl {
		if len(c.FreeIds) > 0 {
			return c.AllocateIntID(), nil
		}
	}
	c, err := d.GetNewChunk()
	if err != nil {
		err := fmt.Errorf("Ids not available")
		return 0, err
	}
	return c.AllocateIntID(), nil
}

func (d *Drsm) ReleaseIntID(sharedPoolName string, id int32) error {
	chunkId := id >> 10
	chunk, found := d.localChunkTbl[chunkId]
	if found == true {
		chunk.ReleaseIntID(id)
		return nil
	}
	err := fmt.Errorf("Unknown Id")
	return err
}

func (d *Drsm) FindOwnerIntID(sharedPoolName string, id int32) (*PodId, error) {
	chunkId := id >> 10
	chunk, found := d.globalChunkTbl[chunkId]
	if found == true {
		podId := chunk.GetOwner()
		return podId, nil
	}
	err := fmt.Errorf("Unknown Id")
	return nil, err
}
