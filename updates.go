package drsm

import (
	"context"
	//"encoding/json"
	"github.com/omec-project/MongoDBLibrary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// handle incoming db notification and update

func handleDbUpdates(d *Drsm) {
	database := MongoDBLibrary.Client.Database(d.db.Name)
	collection := database.Collection(d.sharedPoolName)

	//create stream to monitor actions on the collection
	updateStream, err := collection.Watch(context.TODO(), mongo.Pipeline{bson.D{{"$match",bson.M{{"type", "keepalive"}}}}}) //, options.ChangeStream().SetFullDocument(options.UpdateLookup))

/*
pipeline := mongo.Pipeline{
bson.D{
        {"$match", bson.D{{"$or",
        	bson.A{
        		bson.D{{"fullDocument.username", "alice"}},
        		bson.D{{"operationType", "delete"}}}}},
        }
    }
}
*/

	if err != nil {
		panic(err)
	}
	routineCtx, _ := context.WithCancel(context.Background())
	//run routine to get messages from stream
	iterateChangeStream(d, routineCtx, updateStream)
}

/*
2022/08/24 06:13:58 iterate stream :  map[
                                          _id:map[_data:826305C1A6000000012B022C0100296E5A1004E3379675DE7341FD9485EAD3567681A4463C5F6964003C70756E63684C6976656E657373000004]
                                           clusterTime:{1661321638 1}
                                           documentKey:map[_id:punchLiveness]
                                           fullDocument:map[_id:punchLiveness dbtestapp-7b9c99fcd9-rbw5j:map[podId:dbtestapp-7b9c99fcd9-rbw5j time:1661321638] dbtestapp-7b9c99fcd9-zg98m:map[podId:dbtestapp-7b9c99fcd9-zg98m time:1661321523]]
                                           ns:map[coll:ngapid db:sdcore]
                                           operationType:update
                                           updateDescription:map[
                                                    removedFields:[]
                                                    updatedFields:map[dbtestapp-7b9c99fcd9-rbw5j:map[podId:dbtestapp-7b9c99fcd9-rbw5j time:1661321638]]
                                            ]
                                        ]
*/
func iterateChangeStream(d *Drsm, routineCtx context.Context, stream *mongo.ChangeStream) {
	log.Println("iterate change stream for podData ", d)

	// case 1: Update Global Table in the Drsm with new chunk. There is possiblity of newPod added
	// case 2: If podLivenessCheck results in Pod Down. Then inform Claim Thread.
	// case 2: Update Global Table in the Drsm, existing chunk new owner != ORPHAN i.e. somebody claimed the Chunk
	// case 3: Update Global Table in the Drsm, existing chunk new owner = ORPHAN
	// 	// case 2: Update Global Table in the Drsm, existing chunk new owner != ORPHAN i.e. somebody claimed the Chunk

	defer stream.Close(routineCtx)
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		// If new Pod detected then send it on channel.. d->newPod
		// If existing Pod goes down. d->podDown
		log.Println("iterate stream : ", data)
		for k := range data {
			log.Println("k,v : ", k, data[k])
			log.Println("key matched k,v : ", k, data[k])
			//			for k1 := range data[k] {
			//				log.Println("k1 : ", k1)
			//			}
		}
	}
}

func startDiscovery(d *Drsm) {
	// we discover other endpoints sharing same shared resources through mongoDB streaming
	// Temp : punch liveness in DB every 2 second
	// DB will send notification to every other POD..
	// Every Pod monitors if drsm client missed keepalive for 5 times..
	go punchLiveness(d)
	go checkLiveness(d)
}

func punchLiveness(d *Drsm) {
	// write to DB - signature every 2 second
	ticker := time.NewTicker(2000 * time.Millisecond)
	log.Println("NfProfile document expiry enabled")

	ret := MongoDBLibrary.RestfulAPICreateTTLIndex(d.sharedPoolName, 0, "expireAt")
	if ret {
		log.Println("TTL Index created for Field : expireAt in Collection : NfProfile")
	} else {
		log.Println("TTL Index exists for Field : expireAt in Collection : NfProfile")
	}

	for {
		select {
		case <-ticker.C:
			log.Println("punch liveness goroutine ", d.sharedPoolName)
			filter := bson.M{"_id": d.clientId.PodName}
			now := time.Now()

			timein := time.Now().Local().Add(time.Second * 20)
			t := now.Unix()

			update := bson.D{{"_id", d.clientId.PodName}, {"type", "keepalive"}, {"podId", d.clientId.PodName}, {"time", t}, {"expireAt", timein}}

			_, err := MongoDBLibrary.PutOneCustomDataStructure(d.sharedPoolName, filter, update)
			if err != nil {
				log.Println("put data failed : ", err)
				return
			}
			log.Println("punch liveness goroutine complete")
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
			log.Println("check liveness goroutine ")
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
