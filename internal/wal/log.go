package wal

import (
	"os"
)

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
	// TODO: Implement actual writing to log file
	_ = entry
	return nil
}
