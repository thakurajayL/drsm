package drsm

type Chunk struct {
	Id       int32
	Owner    PodId
	State    ChunkState
	FreeIds  []int32
	AllocIds map[int32]bool
}

func (c *Chunk) GetOwner() {
	return c.Owner

}

func GetNewChunk(d *Drsm) (*Chunk, error) {
	// Get new Chunk
	// We got to allocate new Chunk. We should select
	// probable chunk number

	var id int32 = 10
	c := &Chunk{Id: id}

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
	var i int32 = id & 0x3ff
	c.FreeIds = append(i, c.FreeIds)
}
