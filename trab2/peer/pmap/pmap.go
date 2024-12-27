package pmap

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// StringToTimestampMap represents a map from string to timestamps.
type PMap map[string]time.Time

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
		buffer.WriteString(fmt.Sprintf("  %q: %q,\n", key, value.Format(time.RFC3339)))
	}
	buffer.WriteString("}")
	return buffer.String()
}

// Serialize serializes the map to a byte slice using Gob encoding.
func (m PMap) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize reconstructs a StringToTimestampMap from a byte slice using Gob encoding.
func Deserialize(data []byte) (PMap, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var result PMap
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *PMap) UpdatePeer(peer string) {
    (*m)[peer] = time.Now()
}
