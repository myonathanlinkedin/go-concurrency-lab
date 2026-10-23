package main

import (
    "hash/fnv"
    "sort"
    "strconv"
)

func (ch *ConsistentHash) hashKey(key string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(key))
    return h.Sum64()
}

func (ch *ConsistentHash) AddNode(nodeID string) {
    ch.mutex.Lock()
    defer ch.mutex.Unlock()
    for i := 0; i < ch.virtualNodes; i++ {
        virtualKey := nodeID + "#" + strconv.Itoa(i)
        hash := ch.hashKey(virtualKey)
        ch.ring = append(ch.ring, hash)
        ch.hashMap[hash] = nodeID
    }
    sort.Slice(ch.ring, func(i, j int) bool { return ch.ring[i] < ch.ring[j] })
}

func (ch *ConsistentHash) RemoveNode(nodeID string) {
    ch.mutex.Lock()
    defer ch.mutex.Unlock()
    newRing := make([]uint64, 0, len(ch.ring))
    for _, hash := range ch.ring {
        if ch.hashMap[hash] != nodeID {
            newRing = append(newRing, hash)
        } else {
            delete(ch.hashMap, hash)
        }
    }
    ch.ring = newRing
}

func (ch *ConsistentHash) GetNode(key string) string {
    ch.mutex.RLock()
    defer ch.mutex.RUnlock()
    if len(ch.ring) == 0 {
        return ""
    }
    hash := ch.hashKey(key)
    idx := sort.Search(len(ch.ring), func(i int) bool { return ch.ring[i] >= hash })
    if idx == len(ch.ring) {
        idx = 0
    }
    return ch.hashMap[ch.ring[idx]]
}
