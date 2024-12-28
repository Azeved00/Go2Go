package pmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
    "net"
)

// StringToTimestampMap represents a map from string to timestamps.
type PMap map[net.Addr]time.Time

// Merge combines two StringToTimestampMap structures.
// If the same key exists in both, the latest timestamp is used.
func (m PMap) Merge(other PMap) {
	for key, value := range other {
		if existingValue, exists := m[key]; exists {
			// Keep the later timestamp
			if value.After(existingValue) {
				m[key] = value
			}
		} else {
			m[key] = value
		}
	}
}

// PrettyPrint returns a string representation of the map in a readable format.
func (m PMap) PrettyPrint() string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for key, value := range m {
		buffer.WriteString(fmt.Sprintf("  %q: %q,\n", key, value))
	}
	buffer.WriteString("}")
	return buffer.String()
}

type SerializedPMap struct {
	Address  string    `json:"address"`
	Timestamp string    `json:"timestamp"`
}
func (m PMap) Serialize() ([]byte, error) {
	var serializedEntries []SerializedPMap

	for key, entry := range m {
		serializedEntry := SerializedPMap{
			Address:  key.String(),
			Timestamp: entry.Format(time.RFC3339), 	
        }
		serializedEntries = append(serializedEntries, serializedEntry)
	}

	// Convert the serialized entries to JSON
	return json.Marshal(serializedEntries)
}


func Deserialize(data []byte) (PMap, error) {
	var serializedEntries []SerializedPMap
	err := json.Unmarshal(data, &serializedEntries)
	if err != nil {
		return nil, err
	}

	deserializedMap := make(PMap)
	for _, entry := range serializedEntries {
		addr, err := net.ResolveTCPAddr("tcp", entry.Address) // Assuming TCPAddr for simplicity
		if err != nil {
			return nil, fmt.Errorf("error resolving address: %v", err)
		}

		parsedTime, _ := time.Parse(time.RFC3339, entry.Timestamp)
		deserializedMap[addr] = parsedTime
	}

	return deserializedMap, nil
}

func NewPeerMap() PMap {
    return make(map[net.Addr] time.Time)
}

func (m *PMap) UpdatePeer(peer net.Conn) {
    (*m)[peer.RemoteAddr()] = time.Now()
}
func (m *PMap) DeletePeer(peer net.Conn) {
    defer  delete(*m, peer.RemoteAddr())
}
