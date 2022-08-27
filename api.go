package drsm

import (
	"fmt"
	"log"
)

type DbInfo struct {
	Url  string
	Name string
}

type PodId struct {
	PodName string `bson:"podName,omitempty" json:"podName,omitempty"`
	PodIp   string `bson:"podIp,omitempty" json:"podIp,omitempty"`
}

type DrsmMode int

const (
	ResourceClient DrsmMode = iota + 0
	ResourceDemux
)

type Options struct {
	Mode   DrsmMode
	ScanCb func(int32) bool
}

func InitDRSM(sharedPoolName string, myid PodId, db DbInfo, opt *Options) (*Drsm, error) {
	log.Println("****MY ID ***** ", myid)

	d := &Drsm{sharedPoolName: sharedPoolName,
		clientId: myid,
		db:       db,
		mode:     ResourceClient}

	d.ConstuctDrsm(opt)

	return d, nil
}

func (d *Drsm) AllocateIntID(sharedPoolName string) (int32, error) {
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not allocate Resource index ")
		err := fmt.Errorf("Demux mode does not allow Resource Id allocation")
		return 0, err
	}
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
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not release Resource index ")
		err := fmt.Errorf("Demux mode does not allow Resource Id allocation")
		return err
	}

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
