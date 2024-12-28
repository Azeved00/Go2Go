package pmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type PMap struct { 
   m map[string]time.Time
   Addr string
}

// Merge combines two StringToTimestampMap structures.
// If the same key exists in both, the latest timestamp is used.
func (m PMap) Merge(other PMap) {
	for key, value := range other.m {
		if existingValue, exists := m.m[key]; exists {
			// Keep the later timestamp
			if value.After(existingValue) {
				m.m[key] = value
			}
		} else {
			m.m[key] = value
		}
	}
}

// PrettyPrint returns a string representation of the map in a readable format.
func (m PMap) PrettyPrint() string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for key, value := range m.m {
		buffer.WriteString(fmt.Sprintf("  %q: %q,\n", key, value))
	}
	buffer.WriteString("}")
	return buffer.String()
}

type SerializedPMapEntry struct {
	Address  string    `json:"address"`
	Timestamp string    `json:"timestamp"`
}
type SerializedPMap struct {
	M   []SerializedPMapEntry   `json:"m"`
	Addr string                 `json:"sender"`
}
func (m PMap) Serialize() ([]byte, error) {
	var serializedMap SerializedPMap
    serializedMap.Addr = m.Addr

	for key, entry := range m.m {
		serializedEntry := SerializedPMapEntry{
			Address:  key,
			Timestamp: entry.Format(time.RFC3339), 	
        }
		serializedMap.M = append(serializedMap.M, serializedEntry)
	}

	jsonData, err := json.Marshal(serializedMap)
	if err != nil {
		return nil, err
	}

	// Append a newline character
	jsonData = append(jsonData, '\n')
    return jsonData, nil
}


func Deserialize(data []byte) (PMap, error) {
	var serializedMap SerializedPMap

	err := json.Unmarshal(data, &serializedMap)
	if err != nil {
		return PMap{}, err
	}

	pmap := PMap {
		m: make(map[string]time.Time),
	}

	for _, serializedEntry := range serializedMap.M {

		timestamp, err := time.Parse(time.RFC3339, serializedEntry.Timestamp)
		if err != nil {
			return PMap{}, err
		}

		// Add the entry to the map
		pmap.m[serializedEntry.Address] = timestamp
	}

    pmap.Addr = serializedMap.Addr
	return pmap, nil
}

func NewPeerMap(addr string) PMap {
    return PMap { 
        m: make(map[string] time.Time),
        Addr: addr,
    }
}

func (m *PMap) UpdatePeer(peer string) {
    (*m).m[peer] = time.Now()
}
func (m *PMap) DeletePeer(peer string) {
    defer  delete((*m).m, peer)
}
