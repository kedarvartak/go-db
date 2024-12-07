package storage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// manages how records are serialized and deserialized. converts records coming from the record.go file into bytes and vice versa

// ValueType represents the type of a value in a record
type ValueType byte

const (
	TypeNull ValueType = iota
	TypeInt
	TypeString
	TypeBool
	TypeFloat
)

// SerializedRecord represents a record in its binary form
type SerializedRecord struct {
	data []byte
}

// SerializeRecord converts a record to bytes
func SerializeRecord(record *Record) ([]byte, error) {
	if len(record.Values) == 0 {
		return nil, errors.New("cannot serialize empty record")
	}

	// Calculate required size
	size := 4 // Size for number of values
	typeSection := make([]ValueType, len(record.Values))

	// First pass: calculate size and determine types
	for i, value := range record.Values {
		size++ // Add 1 byte for type information
		switch v := value.(type) {
		case nil:
			typeSection[i] = TypeNull
		case int:
			typeSection[i] = TypeInt
			size += 8
		case string:
			typeSection[i] = TypeString
			size += 4 + len(v) // length prefix + string data
		case bool:
			typeSection[i] = TypeBool
			size += 1
		case float64:
			typeSection[i] = TypeFloat
			size += 8
		default:
			return nil, fmt.Errorf("unsupported type for value: %v", value)
		}
	}

	// Allocate buffer with exact size needed
	buffer := make([]byte, size)
	offset := 0

	// Write number of values
	binary.LittleEndian.PutUint32(buffer[offset:], uint32(len(record.Values)))
	offset += 4

	// Write each value with its type
	for i, value := range record.Values {
		// Write type
		buffer[offset] = byte(typeSection[i])
		offset++

		// Write value based on its type
		switch v := value.(type) {
		case nil:
			// No additional data needed for null
		case int:
			if offset+8 > len(buffer) {
				return nil, errors.New("buffer overflow while writing int")
			}
			binary.LittleEndian.PutUint64(buffer[offset:], uint64(v))
			offset += 8
		case string:
			if offset+4 > len(buffer) {
				return nil, errors.New("buffer overflow while writing string length")
			}
			binary.LittleEndian.PutUint32(buffer[offset:], uint32(len(v)))
			offset += 4
			if offset+len(v) > len(buffer) {
				return nil, errors.New("buffer overflow while writing string data")
			}
			copy(buffer[offset:], v)
			offset += len(v)
		case bool:
			if offset >= len(buffer) {
				return nil, errors.New("buffer overflow while writing bool")
			}
			if v {
				buffer[offset] = 1
			}
			offset++
		case float64:
			if offset+8 > len(buffer) {
				return nil, errors.New("buffer overflow while writing float")
			}
			binary.LittleEndian.PutUint64(buffer[offset:], math.Float64bits(v))
			offset += 8
		}
	}

	return buffer, nil
}

// DeserializeRecord converts bytes back to a record
func DeserializeRecord(data []byte) (*Record, error) {
	if len(data) < 4 {
		return nil, errors.New("invalid record data")
	}

	// Read number of values
	numValues := binary.LittleEndian.Uint32(data)
	offset := 4

	// Read values
	values := make([]interface{}, numValues)
	for i := uint32(0); i < numValues; i++ {
		if offset >= len(data) {
			return nil, errors.New("corrupt record data")
		}

		// Read type
		valueType := ValueType(data[offset])
		offset++

		// Read value
		switch valueType {
		case TypeNull:
			values[i] = nil
		case TypeInt:
			if offset+8 > len(data) {
				return nil, errors.New("corrupt record data")
			}
			values[i] = int(binary.LittleEndian.Uint64(data[offset:]))
			offset += 8
		case TypeString:
			if offset+4 > len(data) {
				return nil, errors.New("corrupt record data")
			}
			length := binary.LittleEndian.Uint32(data[offset:])
			offset += 4
			if offset+int(length) > len(data) {
				return nil, errors.New("corrupt record data")
			}
			values[i] = string(data[offset : offset+int(length)])
			offset += int(length)
		case TypeBool:
			values[i] = data[offset] != 0
			offset++
		case TypeFloat:
			if offset+8 > len(data) {
				return nil, errors.New("corrupt record data")
			}
			bits := binary.LittleEndian.Uint64(data[offset:])
			values[i] = math.Float64frombits(bits)
			offset += 8
		default:
			return nil, errors.New("unknown value type")
		}
	}

	return &Record{Values: values}, nil
}
