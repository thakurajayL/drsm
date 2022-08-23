package drsm
import (
"fmt"
)
func (d *Drsm) podDown() {
	for {
		select {
		case p := <-d.podDown:
			fmt.Println("Pod Down detected ", p)
			// Given Pod find out current Chunks owned by this POD
			pd := d.podMap[p]
			for k, v := range pd.podChunks {
				c, found := d.globalChunkTbl[k]
				fmt.Printf("Found : %v chunk : %v ", found, chunk)
				go c.claimChunk()
			}
		}
	}
}

func (c *chunk) claimChunk(d *Drsm) {
	// try to claim. If success then notification will update owner.
	claimSuccess := true
	if claimSuccess == true {
		d.scanChunk <- c.Id
	}
}
