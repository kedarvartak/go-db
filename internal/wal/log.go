package wal

import (
	"os"
)
// gonna implement this last to the project

// WAL - Write Ahead Logging. Crucial component for database durability and crash recovery.
// log entry records operations that have been committed.
// log type is used to identify the type of operation.
// Before any database change, we first write to WAL
// Each entry has a sequence number for ordering
// If the database crashes, we can replay the WAL to recover
type LogEntry struct {
	Sequence uint64
	Type     LogType
	Data     []byte
}

type WAL struct {
	File     *os.File
	Sequence uint64
}

func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &WAL{
		File:     file,
		Sequence: 0,
	}, nil
}

func (w *WAL) Write(data []byte) error {
	entry := LogEntry{
		Sequence: w.Sequence + 1,
		Type:     LogTypeData,
		Data:     data,
	}
	w.Sequence++
	
	_ = entry
	return nil
}
