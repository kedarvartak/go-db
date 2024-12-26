package wal

// TIRTHRAJ IF YOURE STALKING THIS FUCK YOU GET A LIFE BITCH

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
	"time"
)

// gonna implement this last to the project

// WAL - Write Ahead Logging. Crucial component for database durability and crash recovery.
// log entry records operations that have been committed.
// log type is used to identify the type of operation.
// Before any database change, we first write to WAL
// Each entry has a sequence number for ordering
// If the database crashes, we can replay the WAL to recover
type LogEntry struct {
	LSN       LSN       // Log Sequence Number
	Timestamp time.Time // When the log entry was created
	TxID      uint64    // Transaction ID
	Type      LogType   // Type of log entry
	PageID    uint64    // Affected page ID
	Record    LogRecord // Actual changes
}

type WAL struct {
	mu       sync.Mutex // using mutex to ensure thread safety
	file     *os.File
	filename string
	currentLSN LSN
	buffer   []byte      // Buffer for writing
	bufSize  int        // Size of the buffer
	flushPos int64      // Position of last flush
}
// creates a new WAL instance
func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:     file,
		filename: filename,
		buffer:   make([]byte, 32*1024), // 32KB buffer
		bufSize:  32 * 1024,
	}, nil
}

// For writing log entries to the WAL
func (w *WAL) Write(entry *LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Assign next LSN
	w.currentLSN++
	entry.LSN = w.currentLSN
	entry.Timestamp = time.Now()

	// Serialize the entry
	data, err := w.serializeEntry(entry)
	if err != nil {
		return err
	}

	// Write to file
	if _, err := w.file.Write(data); err != nil {
		return err
	}

	// Force write to disk if this is a commit
	if entry.Type == LogTypeCommitTx {
		return w.file.Sync()
	}

	return nil
}

func (w *WAL) serializeEntry(entry *LogEntry) ([]byte, error) {
	// Calculate total size needed
	size := 8 + // LSN
		8 + // Timestamp
		8 + // TxID
		4 + // Type
		8 + // PageID
		4 + // Record size
		len(entry.Record.Before) +
		len(entry.Record.After)

	buf := make([]byte, size)
	offset := 0

	// Write LSN
	binary.BigEndian.PutUint64(buf[offset:], uint64(entry.LSN))
	offset += 8

	// Write Timestamp
	binary.BigEndian.PutUint64(buf[offset:], uint64(entry.Timestamp.UnixNano()))
	offset += 8

	// Write TxID
	binary.BigEndian.PutUint64(buf[offset:], entry.TxID)
	offset += 8

	// Write Type
	binary.BigEndian.PutUint32(buf[offset:], uint32(entry.Type))
	offset += 4

	// Write PageID
	binary.BigEndian.PutUint64(buf[offset:], entry.PageID)
	offset += 8

	// Write Record
	binary.BigEndian.PutUint32(buf[offset:], uint32(len(entry.Record.Before)))
	offset += 4
	copy(buf[offset:], entry.Record.Before)
	offset += len(entry.Record.Before)

	binary.BigEndian.PutUint32(buf[offset:], uint32(len(entry.Record.After)))
	offset += 4
	copy(buf[offset:], entry.Record.After)

	return buf, nil
}
