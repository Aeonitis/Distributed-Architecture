package server

import (
	"fmt"
	"sync"
)

// Log /**
// A data structure for an append-only sequence of records,
// ordered by time, and you can build a simple commit log with a slice.
type Log struct {
	mu      sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

// Append
// To append a record to the log, you just append to the slice.
// Each time we read a record given an index, we use that index to look up the record in the slice.
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock() // defer function call to the end of the currently executing function & before we return
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

// Read
// Gets record available at offset index
// If the offset given by the client does not exist, return an error
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	currentLengthOfRecords := uint64(len(c.records))
	if offset >= currentLengthOfRecords {
		return Record{}, ErrMessageOffsetNotFound
	}
	return c.records[offset], nil
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrMessageOffsetNotFound = fmt.Errorf("offset not found")
