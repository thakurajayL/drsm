package drsm

import (
	"log"
	"time"
)

func (c *chunk) scanChunk(d *Drsm) {
	if d.mode == ResourceDemux {
		log.Println("Don't perform scan task when demux mode is ON")
		return
	}

	if c.Owner.PodName != d.clientId.PodName {
		log.Println("Don't perform scan task if Chunk is not owned by us")
		return
	}
	c.State = Scanning
	d.scanChunks[c.Id] = c
	var i int32
	for i = 0; i < 1000; i++ {
		c.ScannedIds = append(c.ScannedIds, i)
	}

	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("Lets scan one by one id for %v , chunk details %v ", c.Id, c)
			// TODO : find candidate and then scan that Id.
			// once all Ids are scanned then we can start using this block
			id := c.ScannedIds[len(c.ScannedIds)-1]
			c.ScannedIds = c.ScannedIds[:len(c.ScannedIds)-1]
			rid := c.Id<<10 | id
			res := c.resourceValidCb(rid)
			if res == true {
				c.FreeIds = append(c.FreeIds, id)
			} else {
				c.ScannedIds = append(c.ScannedIds, id) // Moving to the end
			}
		case <-c.stopScan:
			log.Printf("Received Stop Scan. Closing scan for %v", c.Id)
			return
		}
	}
}
