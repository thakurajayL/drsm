# drsm
Distributed Resource Sharing Module (DRSM)

Resources can be
    - integer numbers (TEID, SEID, NGAPIDs,...)
    - IP address pool

Modes
    - demux mode : just listen and get mapping about PODS and their resource assignments
        * can be used by sctplb, upf-adapter
    - client mode : Learn about other clients and their resource mappings
        * can be used by AMF pods, SMF pod

Dependency
    - MongoDB should run in cluster(replicaset) Mode or sharded Mode

Testing
    1. All the DRSM clients discover other clients through pub/sub
    2. Allocate resource id ( indirectly chunk). Other Pods should get notification of newly allocated chunk
    3. POD down event should be detected
    4. Get candidate ORPHAN chunk list once POD down detected
    5. CLAIM chunk to change owner
    6. Through notification other PODS should detect if CHUNK is claimed
    7. Run large number of clients and bring down replicaset by 1..All other pod would try to claim chunks of crashed pod.
       we should see only 1 client claiming it successfully

TODO:
    1. min REST APIs to trigger { allocate, claim }
    2. Rst counter to be appended to identify pod.
    5. Update document needs to figure out if its update for Chunk or update for Keepalive, since we are sharing collection
    7. Clear Separation of demux API vs regular CLIENT API
    8. PodId should be = K8s Pod Id + Rst Count. This makes sure that restarted pod even if it comes with same name then we treat it differently
    9. callback should be available where chunk scanning can be done with help of application
    3. Database module handlign multiple connections.
    6. IP address allocation
    4. PostAPI to accept customData
