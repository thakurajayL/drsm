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

	id := 10
	c := &Chunk{Id: id}

	d.localChunkTbl[id] = c

	// add Ids to freeIds
	return 0
}

func (c *Chunk) AllocateIntID() int32 {
	id := c.FreeIds[len(c.FreeIds)-1]
	c.FreeIds = c.FreeIds[:len(c.FreeIds)-1]
	return (c.Id << 10) | id
}

func (c *Chunk) ReleaseIntID(int id) {
	i := id & 0x3ff
	c.FreeIds = append(i, chunk.FreeIds)
}
