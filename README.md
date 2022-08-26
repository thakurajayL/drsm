# drsm
Distributed Resource sharing go module

Modes
    - demux mode : just listen and get mapping about PODS and their resource assignments
        * can be used by sctplb, upf-adapter
    - client mode : Learn about other clients and their resource mappings
        * can be used by AMF pods, SMF pod

Testing
    1. All the DRSM clients discover other clients through pub/sub
    2. Allocate id ( indirectly chunk). Other Pods should get notification of newly allocated chunk
    3. POD down event should be detected
    4. Get candidate ORPHAN chunk list once POD down detected
    5. CLAIM chunk to change owner
    6. Through notification other PODS should detect if CHUNK is claimed

TODO:
    1. min REST APIs to trigger { allocate, claim }
    2. Rst counter to be appended to identify pod.
    3. Database module handlign multiple connections.
    4. PostAPI to accept customData
