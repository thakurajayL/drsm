package drsm

import (
	"log"
    "time"
)

func (c *Chunk) ScanChunk(d *Drsm) {
	if d.mode == ResourceDemux {
		log.Println("Don't perform scan task when demux mode is ON")
		return
	}

	if c.Owner.PodName != d.clientId.PodName {
		log.Println("Don't perform scan task if Chunk is not owned by us")
		return
	}
	c.State = Scanning

	var maindone chan bool
	var done chan bool

	go func() {
		// keep calling scan callback
		// once done give back signal to main routine
		ticker := time.NewTicker(5000 * time.Millisecond)
		for {
			select {
			case <-maindone:
				log.Printf("Received Stop Scan. Closing scan for %v", c.Id)
			case <-ticker.C:
				log.Printf("Lets scan one by one id for %v", c.Id)
				c.scanCb(c.Id)
				done <- true
			}
		}
	}()

	for {
		select {
		case <-done:
			log.Printf("Scanning complete for %v", c.Id)
			return
		case <-c.stopScan:
			log.Printf("Received Stop Scan. Closing scan for %v", c.Id)
			return
		}
	}
	// add Chunk in the scan list
	// we should plan 1 query/second;
}
