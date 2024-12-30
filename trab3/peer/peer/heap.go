package peer

import (
	//"container/heap"
    "sync"
	//"fmt"
)

type HeapObj struct {
    Counter     uint64
    Word        string
    Addr        string
}

type  ProcessHeap struct {
    heap []HeapObj
    mutex sync.RWMutex
}

func (h *ProcessHeap) Len() int { 
    h.mutex.RLock()
    defer h.mutex.RUnlock()

    return len(h.heap) 
}

func (h *ProcessHeap) Less(i, j int) bool {
    h.mutex.RLock()
    defer h.mutex.RUnlock()

    if h.heap[i].Counter == h.heap[j].Counter {
        return h.heap[i].Word > h.heap[j].Word
    }
    return h.heap[i].Counter > h.heap[j].Counter
}

func (h *ProcessHeap) Swap(i, j int) {
    h.mutex.Lock()
    h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
    h.mutex.Unlock()
}

func (h *ProcessHeap) Push(x HeapObj) {
    h.mutex.Lock()
    h.heap = append(h.heap, x)
    h.mutex.Unlock()
}

func (h *ProcessHeap) Pop() HeapObj {
    h.mutex.Lock()
    old := h.heap
    n := len(old)
    item := old[n-1]
    h.heap = old[0 : n-1]
    h.mutex.Unlock()
    return item
}
