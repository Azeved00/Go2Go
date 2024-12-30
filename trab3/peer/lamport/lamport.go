package lamport

import (
	"bytes"
	"encoding/gob"
	"fmt"
    "sync"
)

type LampClock struct { 
   counter uint64
   lock sync.RWMutex

   Addr string
}


func (m *LampClock) Merge(other *LampClock) {
    m.lock.Lock()
    defer m.lock.Unlock()

    c := other.Get()
    m.counter = max(c, m.counter)
}

// PrettyPrint returns a string representation of the map in a readable format.
func (m *LampClock) PrettyPrint() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return fmt.Sprintf("LampClock [Addr: %s, Counter: %d]", m.Addr, m.counter)
}

func (m *LampClock) Serialize() ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

    var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Serialize only the fields, not the lock itself.
	temp := struct {
		Counter uint64
		Addr    string
	}{
		Counter: m.counter,
		Addr:    m.Addr,
	}

	if err := encoder.Encode(temp); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}


func Deserialize(data []byte) (LampClock, error) {
    var temp struct {
		Counter uint64
		Addr    string
	}

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(&temp); err != nil {
		return LampClock{}, err
	}

	// Initialize and return the LampClock with the deserialized data.
	return LampClock{
		counter: temp.Counter,
		Addr:    temp.Addr,
	}, nil

}

func NewLamportClock(addr string) LampClock{
    return LampClock { 
        Addr: addr,
        counter: 0,
    }
}

func (m *LampClock) Get() uint64 {
    m.lock.RLock()
    defer m.lock.RUnlock()

    return m.counter
}

func (m *LampClock) Increment() uint64 {
    m.lock.Lock()
    defer m.lock.Unlock()

    m.counter = 1 + m.counter
    return m.counter
}
