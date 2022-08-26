package drsm

import (
	"context"
	//"encoding/json"
	"github.com/omec-project/MongoDBLibrary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"strings"
	"time"
)

type UpdatedFields struct {
	ExpireAt time.Time `bson:"expireAt,omitempty"`
	PodId    string    `bson:"podId,omitempty"`
}

type UpdatedDesc struct {
	UpdFields UpdatedFields `bson:"updatedFields,omitempty"`
}

type FullStream struct {
	Id       string    `bson:"_id,omitempty"`
	PodId    string    `bson:"podId,omitempty"`
	ExpireAt time.Time `bson:"expireAt,omitempty"`
	Type     string    `bson:"type,omitempty"`
}

type DocKey struct {
	Id string `bson:"_id,omitempty"`
}

type abc struct {
	DId    DocKey      `bson:"documentKey,omitempty"`
	OpType string      `bson:"operationType,omitempty"`
	Full   FullStream  `bson:"fullDocument,omitempty"`
	Update UpdatedDesc `bson:"updateDescription,omitempty"`
}

/*
 map[
        _id:map[_data:826306F004000000032B022C0100296E5A1004EC0A378B4B3044C28DF4F18548BC3974463C5F6964003C6462746573746170702D6262346334636462342D6A687A6C7A000004]
        clusterTime:{1661399044 3}
        documentKey:map[_id:dbtestapp-bb4c4cdb4-jhzlz]
        ns:map[coll:ngapid db:sdcore]
        operationType:insert
        fullDocument:map[_id:dbtestapp-bb4c4cdb4-jhzlz expireAt:1661399064504 podId:dbtestapp-bb4c4cdb4-jhzlz time:1661399044 type:keepalive]
    ]

map[
        _id:map[_data:826306FE49000000012B022C0100296E5A10045287202787774B43958F3929CFD344D0463C5F6964003C6462746573746170702D3862396634383866372D6337347366000004]
        clusterTime:{1661402697 1}
        documentKey:map[_id:dbtestapp-8b9f488f7-c74sf]
        ns:map[coll:ngapid db:sdcore]
        operationType:update
        updateDescription:map[removedFields:[] updatedFields:map[expireAt:1661402717758 time:1661402697]]
   ]

map[
        _id:map[_data:82630701E5000000012B022C0100296E5A10045287202787774B43958F3929CFD344D0463C5F6964003C6462746573746170702D3862396634383866372D6E64327470000004]
        clusterTime:{1661403621 1}
        documentKey:map[_id:dbtestapp-8b9f488f7-nd2tp]
        ns:map[coll:ngapid db:sdcore]
        operationType:delete
   ]

map[
        _id:map[_data:826307FF400000000B2B022C0100296E5A1004020E4568089B4D8889A42D53E225B5AE463C5F6964003C6368756E6B69642D3131353638000004]
        clusterTime:{1661468480 11}
        documentKey:map[_id:chunkid-11568]
        fullDocument:map[_id:chunkid-11568 podId:dbtestapp-8644b5b7d6-qdk54 type:chunk]
        ns:map[coll:ngapid db:sdcore]
        operationType:insert]


map[
        _id:map[_data:8263085773000000022B022C0100296E5A1004E23062383C624633BDEE5B9B5FEAB2B8463C5F6964003C6368756E6B69642D38333332000004]
        clusterTime:{1661491059 2}
        documentKey:map[_id:chunkid-8332]
        ns:map[coll:ngapid db:sdcore]
        operationType:update
        updateDescription:map[removedFields:[] updatedFields:map[podId:dbtestapp-6dc68f9f68-7fwj8]]]

*/

// handle incoming db notification and update
func (d *Drsm) handleDbUpdates() {
	database := MongoDBLibrary.Client.Database(d.db.Name)
	collection := database.Collection(d.sharedPoolName)

	// TODO : 2 go routines to monitor 2 pipelines
	pipeline := mongo.Pipeline{}

	for {
		//create stream to monitor actions on the collection
		updateStream, err := collection.Watch(context.TODO(), pipeline)

		if err != nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}
		routineCtx, _ := context.WithCancel(context.Background())
		//run routine to get messages from stream
		iterateChangeStream(d, routineCtx, updateStream)
	}
}

func iterateChangeStream(d *Drsm, routineCtx context.Context, stream *mongo.ChangeStream) {
	log.Println("iterate change stream for podData ", d)

	// step 1: update all other clients database
	// case 1: Update Global Table in the Drsm with new chunk. There is possiblity of newPod added
	// case 2: If podLivenessCheck results in Pod Down. Then inform Claim Thread.
	// case 2: Update Global Table in the Drsm, existing chunk new owner != ORPHAN i.e. somebody claimed the Chunk
	// case 3: Update Global Table in the Drsm, existing chunk new owner = ORPHAN
	// case 2: Update Global Table in the Drsm, existing chunk new owner != ORPHAN i.e. somebody claimed the Chunk

	defer stream.Close(routineCtx)
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		var s abc
		bsonBytes, _ := bson.Marshal(data)
		bson.Unmarshal(bsonBytes, &s)
		// If new Pod detected then send it on channel.. d->newPod
		// If existing Pod goes down. d->podDown
		log.Println("iterate stream : ", data)
		log.Printf("\ndecoded stream bson %+v \n", s)
		switch s.OpType {
		case "insert":
			full := &s.Full
			switch full.Type {
			case "keepalive":
				log.Println("insert keepalive document")
				pod, found := d.podMap[full.PodId]
				if found == false {
					podI := PodId{PodName: full.PodId}
					pod := &PodData{PodId: podI}
					pod.podChunks = make(map[int32]*Chunk)
					d.podMap[full.PodId] = pod
					log.Printf("Keepalive insert d.podMaps %+v", d.podMap)
				} else {
					log.Println("insert keepalive document : found existing podId ", pod)
				}
			case "chunk":
				log.Println("insert chunk document")
				pod, found := d.podMap[full.PodId]
				if found == false {
					podI := PodId{PodName: full.PodId}
					pod = &PodData{PodId: podI}
					pod.podChunks = make(map[int32]*Chunk)
					d.podMap[full.PodId] = pod
					log.Printf("Chunk insert  d.podMaps %+v", d.podMap)

				}
				id := full.Id
				log.Println("chunk document id ", id)
				z := strings.Split(id, "-")
				log.Println("extracted chunk id ", z[1])
				cid, _ := strconv.ParseInt(z[1], 10, 32)
				c := int32(cid)
				o := PodId{PodName: full.PodId}
				cp := &Chunk{Id: c, Owner: o}
				pod.podChunks[c] = cp
				d.globalChunkTbl[c] = cp
				log.Println("pod.podChunks ", pod.podChunks)
			}
		case "update":
			log.Println("update operations")
			z := strings.Split(s.DId.Id, "-")
			if len(z) == 2 && z[0] == "chunkid" {
				// update on chunkId..
				// looks like chunk owner getting change
				owner := s.Update.UpdFields.PodId
				cid, _ := strconv.ParseInt(z[1], 10, 32)
				c := int32(cid)
				cp := d.globalChunkTbl[c]
				cp.Owner.PodName = owner // TODO update IP address as well.
				podD := d.podMap[owner]
				podD.podChunks[c] = cp // add chunk to pod
				log.Printf("pod to chunk map %v ", podD.podChunks)
			}
		case "delete":
			log.Println("delete operations")
			x := strings.Split(s.DId.Id, "chunkid")
			if len(x) == 1 {
				// not chunk type doc. So its POD doc.
				// delete olnly gets document id
				pod, found := d.podMap[s.DId.Id]
				log.Printf("Pod %v and  found %v", pod, found)
				log.Printf("Chunks owned by crashed Pod = %v ", pod.podChunks)
				d.podDown <- s.DId.Id
			}
		}
	}
}

func (d *Drsm) startDiscovery() {
	// we discover other endpoints sharing same shared resources through mongoDB streaming
	// Temp : punch liveness in DB every 2 second
	// DB will send notification to every other POD..
	// Every Pod monitors if drsm client missed keepalive for 5 times..
	go punchLiveness(d)
	go checkLiveness(d)
}

func punchLiveness(d *Drsm) {
	// write to DB - signature every 2 second
	ticker := time.NewTicker(20000 * time.Millisecond)

	log.Println(" document expiry enabled")
	ret := MongoDBLibrary.RestfulAPICreateTTLIndex(d.sharedPoolName, 0, "expireAt")
	if ret {
		log.Println("TTL Index created for Field : expireAt in Collection")
	} else {
		log.Println("TTL Index exists for Field : expireAt in Collection")
	}

	for {
		select {
		case <-ticker.C:
			filter := bson.M{"_id": d.clientId.PodName}

			timein := time.Now().Local().Add(time.Second * 20)

			update := bson.D{{"_id", d.clientId.PodName}, {"type", "keepalive"}, {"podId", d.clientId.PodName}, {"expireAt", timein}}

			_, err := MongoDBLibrary.PutOneCustomDataStructure(d.sharedPoolName, filter, update)
			if err != nil {
				log.Println("put data failed : ", err)
				return
			}
		}
	}
}

func checkLiveness(d *Drsm) {
	// go through all pods to see if any pod is showing same old counter
	// Mark it down locally
	// Claiming the chunks can be reactive
	ticker := time.NewTicker(15000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			/*
				for k, p := range d.podMap {
					p.mu.Lock()
					if p.PrevTimestamp.after(p.== p.currentCount {
						d.podDown <- k // let claim thread work chunks assigned by this Pod.
					} else {
						p.prevCount = p.currentCount
					}
					p.mu.Unlock()
				}
			*/
		}
	}
}
