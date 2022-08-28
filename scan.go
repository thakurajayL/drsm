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

	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("Lets scan one by one id for %v", c.Id)
			// TODO : find candidate and then scan that Id.
			// once all Ids are scanned then we can start using this block
			c.scanCb(c.Id)
		case <-c.stopScan:
			log.Printf("Received Stop Scan. Closing scan for %v", c.Id)
			return
		}
	}
}
