package peer

import (
	"encoding/gob"
    "sync"
    "fmt"
)

type LConn struct {
    enc         *gob.Encoder
    lock        sync.RWMutex
}

type ConnMap map[string]*LConn 

func NewConnMap() ConnMap {
    return make(map[string]*LConn)
}

func (m *ConnMap) Len() int{
    return len(*m)
}
func (m *ConnMap) Exists(index string) bool{
    if m == nil { return false }
    x, ok := (*m)[index]
    if !ok { return false }    

    x.lock.Lock()
    defer x.lock.Unlock()
    if x.enc == nil {return false}
    return true
}

func (m *ConnMap) Set(index string, obj *gob.Encoder){
    x, ok := (*m)[index]
    if ok { 
        x.lock.Lock()
        x.enc = obj
        x.lock.Unlock()
    } else {
        (*m)[index] = &(LConn {
            enc : obj,
        })
    } 
}

func (m *ConnMap) Lock(index string) *gob.Encoder{
    x, ok := (*m)[index]
    if !ok {return nil}    

    return x.enc
}

func (m *ConnMap) Unlock(index string) {
    x, ok := (*m)[index]
    if !ok {return}    
    x.lock.Unlock()
}
func (m *ConnMap) Send(index string, message any) bool {
    x, ok := (*m)[index]
    if !ok {return false}    

    x.lock.Lock()
    defer x.lock.Unlock()
    
    if x.enc == nil {
        return false
    }
    if err := x.enc.Encode(message); err != nil {
        fmt.Printf("Failed to encode: %v\n", err)
        return false
    }
    return true
}
func (m *ConnMap) GetKeys() []string{
    var keys []string

    for key, x := range (*m) {
        if x==nil {continue}
        if x.enc==nil {continue}
        keys = append(keys, key)
    }

    return keys
}
